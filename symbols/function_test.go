package symbols_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloakwiss/ntdocs/symbols"
	"github.com/k0kubun/pp/v3"
)

func TestHandleRequriementSectionOfFunction(t *testing.T) {
	fd, er := os.Open("../test/nf-aclapi-treeresetnamedsecurityinfow")
	if er != nil {
		t.Fatal("Cannot open the file")
	}
	defer fd.Close()
	bufFile := bufio.NewReader(fd)

	sections := symbols.GetAllSection(symbols.GetMainContentAsList(bufFile))
	goquery.OuterHtml(sections["requirements"][0])
	table, er := symbols.HandleRequriementSectionOfFunction(sections["requirements"])
	if er != nil {
		t.Fatalf("%s", er.Error())
	}
	mar, er := json.MarshalIndent(table, "", "  ")
	if er != nil {
		t.Fatal("Marshalling failed")
	}
	fmt.Println(string(mar))
}

func TestHandleParameter(t *testing.T) {
	htm := `<div class="content"><p>Waits until one or all of the specified objects are in the signaled state or the time-out interval elapses.</p>
<p>To enter an alertable wait state, use the
<a href="/en-us/windows/desktop/api/synchapi/nf-synchapi-waitformultipleobjectsex" data-linktype="absolute-path">WaitForMultipleObjectsEx</a> function.</p>
<h2 id="syntax">Syntax</h2>
<pre><code class="lang-cpp">DWORD WaitForMultipleObjects(
  [in] DWORD        nCount,
  [in] const HANDLE *lpHandles,
  [in] BOOL         bWaitAll,
  [in] DWORD        dwMilliseconds
);
</code></pre>
<h2 id="parameters">Parameters</h2>
<p><code>[in] nCount</code></p>
<p>The number of object handles in the array pointed to by <i>lpHandles</i>. The maximum number of object handles is <b>MAXIMUM_WAIT_OBJECTS</b>. This parameter cannot be zero.</p>
<p><code>[in] lpHandles</code></p>
<p>An array of object handles. For a list of the object types whose handles can be specified, see the following Remarks section. The array can contain handles to objects of different types. It may not contain multiple copies of the same handle.</p>
<p>If one of these handles is closed while the wait is still pending, the function&#39;s behavior is undefined.</p>
<p>The handles must have the <b>SYNCHRONIZE</b> access right. For more information, see
<a href="/en-us/windows/desktop/SecAuthZ/standard-access-rights" data-linktype="absolute-path">Standard Access Rights</a>.</p>
<p><code>[in] bWaitAll</code></p>
<p>If this parameter is <b>TRUE</b>, the function returns when the state of all objects in the <i>lpHandles</i> array is signaled. If <b>FALSE</b>, the function returns when the state of any one of the objects is set to signaled. In the latter case, the return value indicates the object whose state caused the function to return.</p>
<p><code>[in] dwMilliseconds</code></p>
<p>The time-out interval, in milliseconds. If a nonzero value is specified, the function waits until the specified objects are signaled or the interval elapses. If <i>dwMilliseconds</i> is zero, the function does not enter a wait state if the specified objects are not signaled; it always returns immediately. If <i>dwMilliseconds</i> is <b>INFINITE</b>, the function will return only when the specified objects are signaled.</p>
<p><strong>Windows XP, Windows Server 2003, Windows Vista, Windows 7, Windows Server 2008, and Windows Server 2008 R2:</strong> The <i>dwMilliseconds</i> value does include time spent in low-power states. For example, the timeout does keep counting down while the computer is asleep.</p>
<p><strong>Windows 8 and newer, Windows Server 2012 and newer:</strong> The <i>dwMilliseconds</i> value does not include time spent in low-power states. For example, the timeout does not keep counting down while the computer is asleep.</p>
<h2 id="return-value">Return value</h2>
<p>If the function succeeds, the return value indicates the event that caused the function to return. It can be one of the following values. (Note that <b>WAIT_OBJECT_0</b> is defined as 0 and <b>WAIT_ABANDONED_0</b> is defined as 0x00000080L.)</p>
<table>
<tbody><tr>
<th>Return code/value</th>
<th>Description</th>
</tr>
<tr>
<td width="40%">
<dl>
<dt><b>WAIT_OBJECT_0</b> to (<b>WAIT_OBJECT_0</b> + <i>nCount</i>– 1)</dt>
</dl>
</td>
<td width="60%">
If <i>bWaitAll</i> is <b>TRUE</b>, a return value within the specified range indicates that the state of all specified objects is signaled.
<p>If <i>bWaitAll</i> is <b>FALSE</b>, the return value minus <b>WAIT_OBJECT_0</b> indicates the <i>lpHandles</i> array index of the object that satisfied the wait. If more than one object became signaled during the call, this is the array index of the signaled object with the smallest index value of all the signaled objects.</p>
</td>
</tr>
<tr>
<td width="40%">
<dl>
<dt><b>WAIT_ABANDONED_0</b> to (<b>WAIT_ABANDONED_0</b> + <i>nCount</i>– 1)</dt>
</dl>
</td>
<td width="60%">
If <i>bWaitAll</i> is <b>TRUE</b>, a return value within the specified range indicates that the state of all specified objects is signaled and at least one of the objects is an abandoned mutex object.
<p>If <i>bWaitAll</i> is <b>FALSE</b>, the return value minus <a href="nf-synchapi-waitforsingleobject" data-linktype="relative-path">WAIT_ABANDONED_0</a> indicates the <i>lpHandles</i> array index of an abandoned mutex object that satisfied the wait. Ownership of the mutex object is granted to the calling thread, and the mutex is set to nonsignaled.</p>
<p>If a mutex was protecting persistent state information, you should check it for consistency.</p>
</td>
</tr>
<tr>
<td width="40%">
<dl>
<dt><b>WAIT_TIMEOUT</b></dt>
<dt>0x00000102L</dt>
</dl>
</td>
<td width="60%">
The time-out interval elapsed and the conditions specified by the <i>bWaitAll</i> parameter are not satisfied.
</td>
</tr>
<tr>
<td width="40%">
<dl>
<dt><b>WAIT_FAILED</b></dt>
<dt>(<b>DWORD</b>)0xFFFFFFFF</dt>
</dl>
</td>
<td width="60%">
The function has failed. To get extended error information, call
<a href="/en-us/windows/desktop/api/errhandlingapi/nf-errhandlingapi-getlasterror" data-linktype="absolute-path">GetLastError</a>.
</td>
</tr>
</tbody></table>
<h2 id="remarks">Remarks</h2>
<p>The
<b>WaitForMultipleObjects</b> function determines whether the wait criteria have been met. If the criteria have not been met, the calling thread enters the wait state until the conditions of the wait criteria have been met or the time-out interval elapses.</p>
<p>When <i>bWaitAll</i> is <b>TRUE</b>, the function&#39;s wait operation is completed only when the states of all objects have been set to signaled. The function does not modify the states of the specified objects until the states of all objects have been set to signaled. For example, a mutex can be signaled, but the thread does not get ownership until the states of the other objects are also set to signaled. In the meantime, some other thread may get ownership of the mutex, thereby setting its state to nonsignaled.</p>
<p>When <i>bWaitAll</i> is <b>FALSE</b>, this function checks the handles in the array in order starting with index 0, until one of the objects is signaled. If multiple objects become signaled, the function returns the index of the first handle in the array whose object was signaled.</p>
<p>The function modifies the state of some types of synchronization objects. Modification occurs only for the object or objects whose signaled state caused the function to return. For example, the count of a semaphore object is decreased by one. For more information, see the documentation for the individual synchronization objects.</p>
<p>To wait on more than <b>MAXIMUM_WAIT_OBJECTS</b> handles, use one of the following methods:</p>
<ul>
<li>Create a thread to wait on <b>MAXIMUM_WAIT_OBJECTS</b> handles, then wait on that thread plus the other handles. Use this technique to break the handles into groups of <b>MAXIMUM_WAIT_OBJECTS</b>.</li>
<li>Call <a href="/en-us/windows/desktop/api/winbase/nf-winbase-registerwaitforsingleobject" data-linktype="absolute-path">RegisterWaitForSingleObject</a> or
<a href="/en-us/windows/win32/api/threadpoolapiset/nf-threadpoolapiset-setthreadpoolwait" data-linktype="absolute-path">SetThreadpoolWait</a>
to wait on each handle.
The thread pool waits efficiently on the handles and assigns a worker thread after the object is signaled or the time-out interval expires.</li>
</ul>
The
<b>WaitForMultipleObjects</b> function can specify handles of any of the following object types in the <i>lpHandles</i> array:
<ul>
<li>Change notification</li>
<li>Console input</li>
<li>Event</li>
<li>Memory resource notification</li>
<li>Mutex</li>
<li>Process</li>
<li>Semaphore</li>
<li>Thread</li>
<li>Waitable timer</li>
</ul>
Use caution when calling the wait functions and code that directly or indirectly creates windows. If a thread creates any windows, it must process messages. Message broadcasts are sent to all windows in the system. A thread that uses a wait function with no time-out interval may cause the system to become deadlocked. Two examples of code that indirectly creates windows are DDE and the <a href="/en-us/windows/desktop/api/objbase/nf-objbase-coinitialize" data-linktype="absolute-path">CoInitialize</a> function. Therefore, if you have a thread that creates windows, use
<a href="/en-us/windows/desktop/api/winuser/nf-winuser-msgwaitformultipleobjects" data-linktype="absolute-path">MsgWaitForMultipleObjects</a> or
<a href="/en-us/windows/desktop/api/winuser/nf-winuser-msgwaitformultipleobjectsex" data-linktype="absolute-path">MsgWaitForMultipleObjectsEx</a>, rather than
<b>WaitForMultipleObjects</b>.
<h4 id="examples">Examples</h4>
<p>For an example, see
<a href="/en-us/windows/desktop/Sync/waiting-for-multiple-objects" data-linktype="absolute-path">Waiting for Multiple Objects</a>.</p>
<div class="code"></div>
<h2 id="requirements">Requirements</h2>
<table>
<thead>
<tr>
<th>Requirement</th>
<th style="text-align: left;">Value</th>
</tr>
</thead>
<tbody>
<tr>
<td><strong>Minimum supported client</strong></td>
<td style="text-align: left;">Windows XP [desktop apps | UWP apps]</td>
</tr>
<tr>
<td><strong>Minimum supported server</strong></td>
<td style="text-align: left;">Windows Server 2003 [desktop apps | UWP apps]</td>
</tr>
<tr>
<td><strong>Target Platform</strong></td>
<td style="text-align: left;">Windows</td>
</tr>
<tr>
<td><strong>Header</strong></td>
<td style="text-align: left;">synchapi.h (include Windows.h)</td>
</tr>
<tr>
<td><strong>Library</strong></td>
<td style="text-align: left;">Kernel32.lib</td>
</tr>
<tr>
<td><strong>DLL</strong></td>
<td style="text-align: left;">Kernel32.dll</td>
</tr>
</tbody>
</table>
<h2 id="see-also">See also</h2>
<p><a href="/en-us/windows/win32/api/synchapi/nf-synchapi-waitforsingleobject" data-linktype="absolute-path">WAIT_ABANDONED_0</a></p>
<p><a href="/en-us/windows/desktop/Sync/synchronization-functions" data-linktype="absolute-path">Synchronization Functions</a></p>
<p><a href="/en-us/windows/desktop/Sync/wait-functions" data-linktype="absolute-path">Wait Functions</a></p>
</div>`

	htm1 := `<div class="content"><p>Suspends the specified WOW64 thread.</p>
<h2 id="syntax">Syntax</h2>
<pre><code class="lang-cpp">DWORD Wow64SuspendThread(
  HANDLE hThread
);
</code></pre>
<h2 id="parameters">Parameters</h2>
<p><code>hThread</code></p>
<p>A handle to the thread that is to be suspended. The handle must have the THREAD_SUSPEND_RESUME access right. For more information, see <a href="/en-us/windows/win32/procthread/thread-security-and-access-rights" data-linktype="absolute-path">Thread Security and Access Rights</a>.</p>
<h2 id="return-value">Return value</h2>
<p>If the function succeeds, the return value is the thread&#39;s previous suspend count; otherwise, it is (DWORD) -1. To get extended error information, use the <a href="../errhandlingapi/nf-errhandlingapi-getlasterror" data-linktype="relative-path">GetLastError</a> function.</p>
<h2 id="remarks">Remarks</h2>
<p>If the function succeeds, execution of the specified thread is suspended and the thread&#39;s suspend count is incremented. Suspending a thread causes the thread to stop executing user-mode (application) code.</p>
<p>This function is primarily designed for use by debuggers. It is not intended to be used for thread synchronization. Calling <strong>Wow64SuspendThread</strong> on a thread that owns a synchronization object, such as a mutex or critical section, can lead to a deadlock if the calling thread tries to obtain a synchronization object owned by a suspended thread. To avoid this situation, a thread within an application that is not a debugger should signal the other thread to suspend itself. The target thread must be designed to watch for this signal and respond appropriately.</p>
<p>Each thread has a suspend count (with a maximum value of MAXIMUM_SUSPEND_COUNT). If the suspend count is greater than zero, the thread is suspended; otherwise, the thread is not suspended and is eligible for execution. Calling
<strong>Wow64SuspendThread</strong> causes the target thread&#39;s suspend count to be incremented. Attempting to increment past the maximum suspend count causes an error without incrementing the count.</p>
<p>The <a href="../processthreadsapi/nf-processthreadsapi-resumethread" data-linktype="relative-path">ResumeThread</a> function decrements the suspend count of a suspended thread.</p>
<p>This function is intended for 64-bit applications. It is not supported on 32-bit Windows; such calls fail and set the last error code to ERROR_INVALID_FUNCTION. A 32-bit application can call this function on a WOW64 thread; the result is the same as calling the <a href="../processthreadsapi/nf-processthreadsapi-suspendthread" data-linktype="relative-path">SuspendThread</a> function.</p>
<h2 id="requirements">Requirements</h2>
<table>
<thead>
<tr>
<th>Requirement</th>
<th style="text-align: left;">Value</th>
</tr>
</thead>
<tbody>
<tr>
<td><strong>Minimum supported client</strong></td>
<td style="text-align: left;">Windows Vista</td>
</tr>
<tr>
<td><strong>Minimum supported server</strong></td>
<td style="text-align: left;">Windows Server 2008</td>
</tr>
<tr>
<td><strong>Header</strong></td>
<td style="text-align: left;">wow64apiset.h</td>
</tr>
<tr>
<td><strong>Library</strong></td>
<td style="text-align: left;">Kernel32.lib</td>
</tr>
<tr>
<td><strong>DLL</strong></td>
<td style="text-align: left;">Kernel32.dll</td>
</tr>
</tbody>
</table>
<h2 id="see-also">See also</h2>
<p><a href="../processthreadsapi/nf-processthreadsapi-resumethread" data-linktype="relative-path">ResumeThread</a></p>
</div>`

	_ = htm
	buffer := bufio.NewReader(strings.NewReader(htm1))

	allContent, er := goquery.NewDocumentFromReader(buffer)
	if er != nil {
		log.Panicln("Cannot create the document")
	}
	mainContent := allContent.Find("div.content").First()
	sections := symbols.GetAllSection(symbols.GetContentAsList(mainContent))
	for i, blk := range sections["parameters"] {
		htm, er := blk.Html()
		if er == nil {
			pp.Println(i, htm)
		}
	}
	for _, i := range sections["parameters"][1:2] {
		pp.Println(i.Text())
	}
	if arr, er := symbols.HandleParameterSectionOfFunction(sections["parameters"]); er == nil {
		pp.Println(arr)
	} else {
		log.Println("Comse safwuro")
	}
}
