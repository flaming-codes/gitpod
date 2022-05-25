// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package io.gitpod.jetbrains.remote

import com.intellij.codeWithMe.ClientId
import com.intellij.ide.BrowserUtil
import com.intellij.ide.CommandLineProcessor
import com.intellij.openapi.application.runInEdt
import com.intellij.openapi.client.ClientSession
import com.intellij.openapi.client.ClientSessionsManager
import com.intellij.openapi.diagnostic.thisLogger
import com.intellij.openapi.project.Project
import com.intellij.openapi.util.io.FileUtilRt
import com.intellij.util.application
import io.netty.channel.ChannelHandlerContext
import io.netty.handler.codec.http.FullHttpRequest
import io.netty.handler.codec.http.QueryStringDecoder
import kotlinx.coroutines.DelicateCoroutinesApi
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import org.jetbrains.ide.RestService
import java.nio.file.InvalidPathException
import java.nio.file.Path

@Suppress("OPT_IN_IS_NOT_ENABLED", "UnstableApiUsage")
@OptIn(DelicateCoroutinesApi::class)
class GitpodCLIService : RestService() {

    override fun getServiceName() = SERVICE_NAME

    override fun execute(urlDecoder: QueryStringDecoder, request: FullHttpRequest, context: ChannelHandlerContext): String? {
        if (application.isHeadlessEnvironment) {
            return "not supported in headless mode"
        }
        val operation = getStringParameter("op", urlDecoder)
        if (operation == "open") {
            val fileStr = getStringParameter("file", urlDecoder)
            if (fileStr.isNullOrBlank()) {
                return "file is missing"
            }
            val file = parseFilePath(fileStr) ?: return "invalid file"
            val shouldWait = getBooleanParameter("wait", urlDecoder)
            return withClient(request, context) {
                CommandLineProcessor.doOpenFileOrProject(file, shouldWait).future.get()
            }
        }
        if (operation == "preview") {
            val url = getStringParameter("url", urlDecoder)
            if (url.isNullOrBlank()) {
                return "url is missing"
            }
            return withClient(request, context) { project ->
                BrowserUtil.browse(url, project)
            }
        }
        return "invalid operation"
    }

    private tailrec suspend fun getClientSessionAsync(): ClientSession {
        val project = getLastFocusedOrOpenedProject()
        var session: ClientSession? = null
        if (project != null) {
            session = ClientSessionsManager.getProjectSessions(project, false).firstOrNull()
        }
        if (session == null) {
            session = ClientSessionsManager.getAppSessions(false).firstOrNull()
        }
        return if (session != null) {
            session
        } else {
            delay(1000L)
            getClientSessionAsync()
        }
    }

    private fun withClient(request: FullHttpRequest, context: ChannelHandlerContext, action: (project: Project?) -> Unit): String? {
        GlobalScope.launch {
            val session = getClientSessionAsync()
            runInEdt {
                ClientId.withClientId(session.clientId) {
                    action(getLastFocusedOrOpenedProject())
                    sendOk(request, context)
                }
            }
        }
        return null
    }

    private fun parseFilePath(path: String): Path? {
        return try {
            var file: Path = Path.of(FileUtilRt.toSystemDependentName(path)) // handle paths like '/file/foo\qwe'
            if (!file.isAbsolute) {
                file = file.toAbsolutePath()
            }
            file.normalize()
        } catch (e: InvalidPathException) {
            thisLogger().warn("gitpod cli: failed to parse file path:", e)
            null
        }
    }

    companion object {
        const val SERVICE_NAME = "gitpod/cli"
    }
}
