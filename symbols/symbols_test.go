package symbols_test

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/cloakwiss/ntdocs/symbols"
)

const htm string = `<div class="content"><p>Updates the specified attribute in a list of attributes for process and thread creation.</p>
<h2 id="syntax">Syntax</h2>
<pre><code class="lang-cpp">BOOL UpdateProcThreadAttribute(
 [in, out]       LPPROC_THREAD_ATTRIBUTE_LIST lpAttributeList,
 [in]            DWORD                        dwFlags,
 [in]            DWORD_PTR                    Attribute,
 [in]            PVOID                        lpValue,
 [in]            SIZE_T                       cbSize,
 [out, optional] PVOID                        lpPreviousValue,
 [in, optional]  PSIZE_T                      lpReturnSize
);
</code></pre>
<h2 id="parameters">Parameters</h2>
<p><code>[in, out] lpAttributeList</code></p>
<p>A pointer to an attribute list created by the <a href="/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-initializeprocthreadattributelist" data-linktype="absolute-path">InitializeProcThreadAttributeList</a> function.</p>
<p><code>[in] dwFlags</code></p>
<p>This parameter is reserved and must be zero.</p>
<p><code>[in] Attribute</code></p>
<p>The attribute key to update in the attribute list. This parameter can be one of the following values.</p>
<table>
<tbody><tr>
<th>Value</th>
<th>Meaning</th>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_GROUP_AFFINITY"></a><a id="proc_thread_attribute_group_affinity"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_GROUP_AFFINITY</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a <a href="/en-us/windows/desktop/api/winnt/ns-winnt-group_affinity" data-linktype="absolute-path">GROUP_AFFINITY</a> structure that specifies the processor group affinity for the new thread.
<p>Supported in Windows 7 and newer and Windows Server 2008 R2 and newer.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_HANDLE_LIST"></a><a id="proc_thread_attribute_handle_list"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_HANDLE_LIST</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a list of handles to be inherited by the child process.
<p>These handles must be created as inheritable handles and must not include pseudo handles such as those returned by the <a href="/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-getcurrentprocess" data-linktype="absolute-path">GetCurrentProcess</a> or <a href="/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-getcurrentthread" data-linktype="absolute-path">GetCurrentThread</a> function.</p>
<div class="alert"><b>Note</b>  if you use this attribute, pass in a value of TRUE for the <i>bInheritHandles</i> parameter of the <a href="/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-createprocessa" data-linktype="absolute-path">CreateProcess</a> function.</div>
<div> </div>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_IDEAL_PROCESSOR"></a><a id="proc_thread_attribute_ideal_processor"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_IDEAL_PROCESSOR</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a  <a href="/en-us/windows/desktop/api/winnt/ns-winnt-processor_number" data-linktype="absolute-path">PROCESSOR_NUMBER</a> structure that specifies the ideal processor for the new thread.
<p>Supported in Windows 7 and newer and Windows Server 2008 R2 and newer.</p>
</td>
</tr>

<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_MACHINE_TYPE"></a><a id="proc_thread_attribute_machine_type"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_MACHINE_TYPE</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a <b>WORD</b> that specifies the machine architecture of the child process.
<p>Supported in Windows 11 and newer.</p>
<p>The  <b>WORD</b> pointed to by <i>lpValue</i> can be a value listed on <a href="/en-us/windows/win32/sysinfo/image-file-machine-constants" data-linktype="absolute-path">IMAGE FILE MACHINE CONSTANTS</a>.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_MITIGATION_POLICY"></a><a id="proc_thread_attribute_mitigation_policy"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_MITIGATION_POLICY</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a <b>DWORD</b> or <b>DWORD64</b> that specifies the exploit mitigation policy for the child process. Starting in Windows 10, version 1703, this parameter can also be a pointer to a two-element <b>DWORD64</b> array.
<p>The specified policy overrides the policies set for the application and the system and cannot be changed after the child process starts running.</p>
<p>The  <b>DWORD</b> or <b>DWORD64</b> pointed to by <i>lpValue</i> can be one or more of the values listed in the remarks.</p>
<p>Supported in Windows 7 and newer and Windows Server 2008 R2 and newer.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_PARENT_PROCESS"></a><a id="proc_thread_attribute_parent_process"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_PARENT_PROCESS</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to the handle of a process to use (instead of the calling process) as the parent for the process being created. The handle for the process used must have the <b>PROCESS_CREATE_PROCESS</b> access right.
<p>Attributes inherited from the specified process include handles, the device map, processor affinity, priority, quotas, the process token, and job object. (Note that some attributes such as the debug port will come from the creating process, not the process specified by this handle.)</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_PREFERRED_NODE"></a><a id="proc_thread_attribute_preferred_node"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_PREFERRED_NODE</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to the node number of the preferred NUMA node for the new process.
<p>Supported in Windows 7 and newer and Windows Server 2008 R2 and newer.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_UMS_THREAD"></a><a id="proc_thread_attribute_ums_thread"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_UMS_THREAD</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a <a href="/en-us/windows/desktop/api/winnt/ns-winnt-ums_create_thread_attributes" data-linktype="absolute-path">UMS_CREATE_THREAD_ATTRIBUTES</a> structure that specifies a user-mode scheduling (UMS) thread context and a UMS completion list to associate with the thread.
<p>After the UMS thread is created, the system queues it to the specified completion list. The UMS thread runs only when an application&#39;s UMS scheduler retrieves the UMS thread from the completion list and selects it to run.  For more information, see <a href="/en-us/windows/desktop/ProcThread/user-mode-scheduling" data-linktype="absolute-path">User-Mode Scheduling</a>.</p>
<p>Supported in Windows 7 and newer and Windows Server 2008 R2 and newer.</p>
<p>Not supported in Windows 11 and newer (see <a href="/en-us/windows/win32/procthread/user-mode-scheduling" data-linktype="absolute-path">User-Mode Scheduling</a>).</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_SECURITY_CAPABILITIES"></a><a id="proc_thread_attribute_security_capabilities"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_SECURITY_CAPABILITIES</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a <a href="/en-us/windows/desktop/api/winnt/ns-winnt-security_capabilities" data-linktype="absolute-path">SECURITY_CAPABILITIES</a> structure that defines the security capabilities of an app container. If this attribute is set the new process will be created as an AppContainer process.
<p>Supported in Windows 8 and newer and Windows Server 2012 and newer.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_PROTECTION_LEVEL"></a><a id="proc_thread_attribute_protection_level"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_PROTECTION_LEVEL</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a <b>DWORD</b> value of <b>PROTECTION_LEVEL_SAME</b>. This specifies the protection level of the child process to be the same as the protection level of its parent process.
<p>Supported in Windows 8.1 and newer and Windows Server 2012 R2 and newer.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_CHILD_PROCESS_POLICY"></a><a id="proc_thread_attribute_child_process_policy"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_CHILD_PROCESS_POLICY</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a <b>DWORD</b> value that specifies the child process policy. The policy specifies whether to allow a child process to be created.
<p>For information on the possible values for the <b>DWORD</b> to which <i>lpValue</i> points, see Remarks.</p>
<p>Supported in Windows 10 and newer and Windows Server 2016 and newer.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_DESKTOP_APP_POLICY"></a><a id="proc_thread_attribute_desktop_app_policy"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_DESKTOP_APP_POLICY</b></dt>
</dl>
</td>
<td width="60%">
This attribute is relevant only to win32 applications that have been converted to UWP packages by using the <a href="https://developer.microsoft.com/windows/bridges/desktop" data-linktype="external">Desktop Bridge</a>.
<p>The <i>lpValue</i> parameter is a pointer to a <b>DWORD</b> value that specifies the desktop app policy. The policy specifies whether descendant processes should continue to run in the desktop environment.</p>
<p>For information about the possible values for the <b>DWORD</b> to which <i>lpValue</i> points, see Remarks.</p>
<p>Supported in Windows 10 Version 1703 and newer and Windows Server Version 1709 and newer.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_JOB_LIST"></a><a id="proc_thread_attribute_job_list"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_JOB_LIST</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a list of job handles to be assigned to the child process, in the order specified.
<p>Supported in Windows 10 and newer and Windows Server 2016 and newer.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="PROC_THREAD_ATTRIBUTE_ENABLE_OPTIONAL_XSTATE_FEATURES"></a><a id="proc_thread_attribute_enable_optional_xstate_features"></a><dl>
<dt><b>PROC_THREAD_ATTRIBUTE_ENABLE_OPTIONAL_XSTATE_FEATURES</b></dt>
</dl>
</td>
<td width="60%">
The <i>lpValue</i> parameter is a pointer to a <b>DWORD64</b> value that specifies the set of optional XState features to enable for the new thread.
<p>Supported in Windows 11 and newer and Windows Server 2022 and newer.</p>
</td>
</tr>
</tbody></table>
<p><code>[in] lpValue</code></p>
<p>A pointer to the attribute value. <b>This value must persist until the attribute list is destroyed using the <a href="/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-deleteprocthreadattributelist" data-linktype="absolute-path">DeleteProcThreadAttributeList</a> function</b>.</p>
<p><code>[in] cbSize</code></p>
<p>The size of the attribute value specified by the <i>lpValue</i> parameter.</p>
<p><code>[out, optional] lpPreviousValue</code></p>
<p>This parameter is reserved and must be NULL.</p>
<p><code>[in, optional] lpReturnSize</code></p>
<p>This parameter is reserved and must be NULL.</p>
<h2 id="return-value">Return value</h2>
<p>If the function succeeds, the return value is nonzero.</p>
<p>If the function fails, the return value is zero. To get extended error information, call
<a href="/en-us/windows/desktop/api/errhandlingapi/nf-errhandlingapi-getlasterror" data-linktype="absolute-path">GetLastError</a>.</p>
<h2 id="remarks">Remarks</h2>
<p>An attribute list is an opaque structure that consists of a series of key/value pairs, one for each attribute. A process can update only the attribute keys described in this topic.</p>
<p>The  <b>DWORD</b> or <b>DWORD64</b> pointed to by <i>lpValue</i> can be one or more of the following values when you specify <b>PROC_THREAD_ATTRIBUTE_MITIGATION_POLICY</b> for the <i>Attribute</i> parameter:</p><dl><p></p>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_DEP_ENABLE</b> (0x00000001)Enables data execution prevention (DEP) for the child process. For more information, see <a href="/en-us/windows/desktop/Memory/data-execution-prevention" data-linktype="absolute-path">Data Execution Prevention</a>.
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_DEP_ATL_THUNK_ENABLE</b> (0x00000002)Enables DEP-ATL thunk emulation for the child process. DEP-ATL thunk emulation causes the system to intercept NX faults that originate from the Active Template Library (ATL) thunk layer. This value can be specified only with PROCESS_CREATION_MITIGATION_POLICY_DEP_ENABLE.
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_SEHOP_ENABLE</b> (0x00000004)Enables structured exception handler overwrite protection (SEHOP) for the child process. SEHOP blocks exploits that use the structured exception handler (SEH) overwrite technique.
</dd>
</dl>
<dl>
<dd>
<b>Windows 7, Windows Server 2008 R2, Windows Server 2008 and Windows Vista:  </b>The following values are not supported until  Windows 8 and Windows Server 2012.
<dl>
<dd>
The force Address Space Layout Randomization (ASLR) policy, if enabled, forcibly rebases images that  are not dynamic base compatible by acting as though an image base  collision happened at load time.  If relocations are required, images that do not have  a base relocation section will not be loaded.
<p>The following mitigation options are available for mandatory ASLR policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_FORCE_RELOCATE_IMAGES_ALWAYS_ON</b> (0x00000001 &lt;&lt;  8)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_FORCE_RELOCATE_IMAGES_ALWAYS_OFF</b> (0x00000002 &lt;&lt;  8)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_FORCE_RELOCATE_IMAGES_ALWAYS_ON_REQ_RELOCS</b> (0x00000003 &lt;&lt;  8)</dd>
</dl>
</dd>
<dd>
The heap terminate on corruption policy, if enabled, causes the heap to terminate if it becomes corrupt.  Note that &#39;always off&#39; does  not override the default opt-in for binaries with current subsystem versions  set in the image header.  Heap terminate on corruption is user mode enforced.
<p>The following mitigation options are available for heap terminate on corruption policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_HEAP_TERMINATE_ALWAYS_ON</b> (0x00000001 &lt;&lt; 12)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_HEAP_TERMINATE_ALWAYS_OFF</b> (0x00000002 &lt;&lt; 12)</dd>
</dl>
</dd>
<dd>
The bottom-up randomization policy, which includes stack randomization options,  causes a random location to be used as the lowest user address.
<p>The following mitigation options are available for the bottom-up randomization policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_BOTTOM_UP_ASLR_ALWAYS_ON</b> (0x00000001 &lt;&lt; 16)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_BOTTOM_UP_ASLR_ALWAYS_OFF</b> (0x00000002 &lt;&lt; 16)</dd>
</dl>
</dd>
<dd>
The high-entropy bottom-up randomization policy, if enabled, causes up to 1TB of bottom-up variance to be used.  Note that high-entropy bottom-up randomization is effective if and only if bottom-up ASLR is also enabled; high-entropy bottom-up randomization is only meaningful for native 64-bit processes.
<p>The following mitigation options are available for the high-entropy bottom-up randomization policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_HIGH_ENTROPY_ASLR_ALWAYS_ON</b> (0x00000001 &lt;&lt; 20)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_HIGH_ENTROPY_ASLR_ALWAYS_OFF                   </b> (0x00000002 &lt;&lt; 20)</dd>
</dl>
</dd>
<dd>
The strict handle checking enforcement policy, if enabled, causes an exception to be raised immediately on a bad handle reference. If this policy is not enabled, a failure status will be returned from the handle reference instead.
<p>The following mitigation options are available for the strict handle checking enforcement policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_STRICT_HANDLE_CHECKS_ALWAYS_ON</b> (0x00000001 &lt;&lt; 24)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_STRICT_HANDLE_CHECKS_ALWAYS_OFF</b> (0x00000002 &lt;&lt; 24)</dd>
</dl>
</dd>
<dd>
The Win32k system call disable policy, if enabled, prevents a process from making Win32k calls.
<p>The following mitigation options are available for the Win32k system call disable policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_WIN32K_SYSTEM_CALL_DISABLE_ALWAYS_ON</b> (0x00000001 &lt;&lt; 28)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_WIN32K_SYSTEM_CALL_DISABLE_ALWAYS_OFF</b> (0x00000002 &lt;&lt; 28)</dd>
</dl>
</dd>
<dd>
The Extension Point Disable policy, if enabled, prevents certain built-in third party extension points from being used.  This policy blocks the following extension points:
<ul>
<li>AppInit DLLs</li>
<li>Winsock Layered Service Providers (LSPs)</li>
<li>Global Windows Hooks</li>
<li>Legacy Input Method Editors (IMEs)</li>
</ul>
Local hooks still work with the Extension Point Disable policy enabled. This behavior is used to prevent legacy extension points from being loaded into a process that does not use them.
<p>The following mitigation options are available for the extension point disable policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_EXTENSION_POINT_DISABLE_ALWAYS_ON</b> (0x00000001 &lt;&lt; 32) </dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_EXTENSION_POINT_DISABLE_ALWAYS_OFF</b> (0x00000002 &lt;&lt; 32) </dd>
</dl>
</dd>
<dd>
The <a href="/en-us/windows/desktop/SecBP/control-flow-guard" data-linktype="absolute-path">Control Flow Guard (CFG) policy</a>, if turned on, places additional restrictions on indirect calls in code that has been built with CFG enabled.
<p>The following mitigation options are available for controlling the CFG policy:</p>
<ul>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_CONTROL_FLOW_GUARD_MASK</b> (0x00000003ui64 &lt;&lt; 40)</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_CONTROL_FLOW_GUARD_DEFER</b> (0x00000000ui64 &lt;&lt; 40)
</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_CONTROL_FLOW_GUARD_ALWAYS_ON</b> (0x00000001ui64 &lt;&lt; 40)</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_CONTROL_FLOW_GUARD_ALWAYS_OFF</b> (0x00000002ui64 &lt;&lt; 40)
</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_CONTROL_FLOW_GUARD_EXPORT_SUPPRESSION</b> (0x00000003ui64 &lt;&lt; 40)
</li>
</ul>
</dd>
<dd>
In addition, the following policy can be specified to enforce that EXEs/DLLs must enable CFG. If an attempt is made to load an EXE/DLL that does not enable CFG, the load will fail:
<ul>
<li><b>PROCESS_CREATION_MITIGATION_POLICY2_STRICT_CONTROL_FLOW_GUARD_MASK</b> (0x00000003ui64 &lt;&lt; 8)</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY2_STRICT_CONTROL_FLOW_GUARD_DEFER</b> (0x00000000ui64 &lt;&lt; 8)
</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY2_STRICT_CONTROL_FLOW_GUARD_ALWAYS_ON</b> (0x00000001ui64 &lt;&lt; 8)</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY2_STRICT_CONTROL_FLOW_GUARD_ALWAYS_OFF</b> (0x00000002ui64 &lt;&lt; 8)
</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY2_STRICT_CONTROL_FLOW_GUARD_RESERVED</b> (0x00000003ui64 &lt;&lt; 8)
</li>
</ul>
</dd>
<dd>
The dynamic code policy, if turned on, prevents a process from generating dynamic code or modifying executable code.
<p>The following mitigation options are available for the dynamic code policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_PROHIBIT_DYNAMIC_CODE_MASK</b> (0x00000003ui64 &lt;&lt; 36)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_PROHIBIT_DYNAMIC_CODE_DEFER</b> (0x00000000ui64 &lt;&lt; 36)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_PROHIBIT_DYNAMIC_CODE_ALWAYS_ON</b> (0x00000001ui64 &lt;&lt; 36)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_PROHIBIT_DYNAMIC_CODE_ALWAYS_OFF</b> (0x00000002ui64 &lt;&lt; 36)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_PROHIBIT_DYNAMIC_CODE_ALWAYS_ON_ALLOW_OPT_OUT</b> (0x00000003ui64 &lt;&lt; 36)
</dd>
</dl>
</dd>
<dd>
The binary signature policy requires EXEs/DLLs to be properly signed.
<p>The following mitigation options are available for the binary signature policy:</p>
<ul>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_MASK</b> (0x00000003ui64 &lt;&lt; 44)</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_DEFER</b> (0x00000000ui64 &lt;&lt; 44)
</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_ALWAYS_ON</b> (0x00000001ui64 &lt;&lt; 44)</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_ALWAYS_OFF</b> (0x00000002ui64 &lt;&lt; 44)
</li>
<li><b>PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_ALLOW_STORE</b> (0x00000003ui64 &lt;&lt; 44)
</li>
</ul>
</dd>
<dd>
The font loading prevention policy for the process determines whether non-system fonts can be loaded for a process. When the policy is turned on, the process is prevented from loading non-system fonts.
<p>The following mitigation options are available for the font loading prevention policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_FONT_DISABLE_MASK</b>                              (0x00000003ui64 &lt;&lt; 48)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_FONT_DISABLE_DEFER</b>                             (0x00000000ui64 &lt;&lt; 48)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_FONT_DISABLE_ALWAYS_ON</b>                         (0x00000001ui64 &lt;&lt; 48)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_FONT_DISABLE_ALWAYS_OFF</b>                        (0x00000002ui64 &lt;&lt; 48)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_AUDIT_NONSYSTEM_FONTS</b>                          (0x00000003ui64 &lt;&lt; 48)
</dd>
</dl>
</dd>
<dd>
The image loading policy of the process determines the types of executable images that can be mapped into the process. When the policy is turned on, images cannot be loaded from some locations, such as remove devices or files that have the Low mandatory label.
<p>The following mitigation options are available for the image loading policy:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_REMOTE_MASK</b>                      (0x00000003ui64 &lt;&lt; 52)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_REMOTE_DEFER</b>                     (0x00000000ui64 &lt;&lt; 52)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_REMOTE_ALWAYS_ON</b>                 (0x00000001ui64 &lt;&lt; 52)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_REMOTE_ALWAYS_OFF</b>                (0x00000002ui64 &lt;&lt; 52)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_REMOTE_RESERVED</b>                  (0x00000003ui64 &lt;&lt; 52)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_LOW_LABEL_MASK</b>                   (0x00000003ui64 &lt;&lt; 56)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_LOW_LABEL_DEFER</b>                  (0x00000000ui64 &lt;&lt; 56)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_LOW_LABEL_ALWAYS_ON              </b>(0x00000001ui64 &lt;&lt; 56)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_LOW_LABEL_ALWAYS_OFF</b>             (0x00000002ui64 &lt;&lt; 56)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_NO_LOW_LABEL_RESERVED               </b>(0x00000003ui64 &lt;&lt; 56)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_PREFER_SYSTEM32_MASK</b>                (0x00000003ui64 &lt;&lt; 60)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_PREFER_SYSTEM32_DEFER               </b>(0x00000000ui64 &lt;&lt; 60)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_PREFER_SYSTEM32_ALWAYS_ON</b>           (0x00000001ui64 &lt;&lt; 60)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_PREFER_SYSTEM32_ALWAYS_OFF</b>          (0x00000002ui64 &lt;&lt; 60)
</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY_IMAGE_LOAD_PREFER_SYSTEM32_RESERVED</b>            (0x00000003ui64 &lt;&lt; 60)
</dd>
</dl>
</dd>
</dl>
</dd>
<dd>
<b>Windows 10, version 1709:  </b>The following value is available only in  Windows 10, version 1709 or later and only with the January 2018 Windows security updates and any applicable firmware updates from the OEM device manufacturer. See <a href="https://support.microsoft.com/help/4073119/protect-against-speculative-execution-side-channel-vulnerabilities-in" data-linktype="external">Windows Client Guidance for IT Pros to protect against speculative execution side-channel vulnerabilities</a>.
<dl>
<dd>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_RESTRICT_INDIRECT_BRANCH_PREDICTION_ALWAYS_ON </b> (0x00000001ui64 &lt;&lt; 16)This flag can be used by processes to protect against sibling hardware threads (hyperthreads) from interfering with indirect branch predictions. Processes that have sensitive information in their address space should consider enabling this flag to protect against attacks involving indirect branch prediction (such as CVE-2017-5715).
</dd>
</dl>
</dd>
</dl>
</dd>
<dd>
<b>Windows 10, version 1809:  </b>The following value is available only in  Windows 10, version 1809 or later.
<dl>
<dd>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_SPECULATIVE_STORE_BYPASS_DISABLE_ALWAYS_ON </b> (0x00000001ui64 &lt;&lt; 24)This flag can be used by processes to disable the Speculative Store Bypass (SSB) feature of CPUs that may be vulnerable to speculative execution side channel attacks involving SSB (CVE-2018-3639). This flag is only supported by certain Intel CPUs that have the requisite hardware features. On CPUs that do not support this feature, the flag has no effect.
</dd>
</dl>
</dd>
</dl>
</dd>
</dl>
<p><b>Windows 10, version 2004:  </b>The following values are available only in  Windows 10, version 2004 or later.</p>
<p>Hardware-enforced Stack Protection (HSP) is a hardware-based security feature where the CPU verifies function return addresses at runtime by employing a shadow stack mechanism.
For user-mode HSP, the default mode is compatibility mode, where only shadow stack violations occurring in modules that are considered compatible with shadow stacks (CETCOMPAT) are fatal.
In strict mode, all shadow stack violations are fatal.</p>
<p>The following mitigation options are available for user-mode Hardware-enforced Stack Protection and related features:</p>
<dl>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_CET_USER_SHADOW_STACKS_ALWAYS_ON </b> (0x00000001ui64 &lt;&lt; 28)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_CET_USER_SHADOW_STACKS_ALWAYS_OFF </b> (0x00000002ui64 &lt;&lt; 28)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_CET_USER_SHADOW_STACKS_STRICT_MODE </b> (0x00000003ui64 &lt;&lt; 28)</dd>
<p>Instruction Pointer validation:</p>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_USER_CET_SET_CONTEXT_IP_VALIDATION_ALWAYS_ON </b>     (0x00000001ui64 &lt;&lt; 32)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_USER_CET_SET_CONTEXT_IP_VALIDATION_ALWAYS_OFF </b>    (0x00000002ui64 &lt;&lt; 32)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_USER_CET_SET_CONTEXT_IP_VALIDATION_RELAXED_MODE </b>  (0x00000003ui64 &lt;&lt; 32)</dd>
<p>Blocking the load of non-CETCOMPAT/non-EHCONT binaries:</p>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_BLOCK_NON_CET_BINARIES_ALWAYS_ON  </b>               (0x00000001ui64 &lt;&lt; 36)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_BLOCK_NON_CET_BINARIES_ALWAYS_OFF </b>               (0x00000002ui64 &lt;&lt; 36)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_BLOCK_NON_CET_BINARIES_NON_EHCONT </b>               (0x00000003ui64 &lt;&lt; 36)</dd>
<p>Restricting certain HSP APIs used to specify security properties of dynamic code to only be callable from outside of the process:</p>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_CET_DYNAMIC_APIS_OUT_OF_PROC_ONLY_ALWAYS_ON </b>     (0x00000001ui64 &lt;&lt; 48)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_CET_DYNAMIC_APIS_OUT_OF_PROC_ONLY_ALWAYS_OFF </b>    (0x00000002ui64 &lt;&lt; 48)</dd>
<p>The FSCTL system call disable policy, if enabled, prevents a process from making NtFsControlFile calls.
The following mitigation options are available for the FSCTL system call disable policy:</p>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_FSCTL_SYSTEM_CALL_DISABLE_ALWAYS_ON </b> (0x00000001ui64 &lt;&lt; 56)</dd>
<dd><b>PROCESS_CREATION_MITIGATION_POLICY2_FSCTL_SYSTEM_CALL_DISABLE_ALWAYS_OFF </b> (0x00000002ui64 &lt;&lt; 56)</dd>
<p>The <b>DWORD</b> pointed to by <i>lpValue</i> can be one or more of the following values when you specify <b>PROC_THREAD_ATTRIBUTE_CHILD_PROCESS_POLICY</b> for the <i>Attribute</i> parameter:</p>
<p><b>PROCESS_CREATION_CHILD_PROCESS_RESTRICTED</b>                                         0x01</p>
<p>The process being created is not allowed to create child processes.  This restriction becomes a property of the token as which the process runs. It should be noted that this restriction is only effective in sandboxed applications (such as AppContainer) which ensure privileged process handles are not accessible to the process. For example, if a process restricting child process creation is able to access another process handle with PROCESS_CREATE_PROCESS or PROCESS_VM_WRITE access rights, then it may be possible to bypass the child process restriction.</p>
<p><b>PROCESS_CREATION_CHILD_PROCESS_OVERRIDE</b>                                           0x02</p>
<p>The process being created is allowed to create a child process, if it would otherwise be restricted. You can only specify this value if the process that is creating the new process is not restricted.</p>
<p>The  <b>DWORD</b> pointed to by <i>lpValue</i> can be one or more of the following values when you specify <b>PROC_THREAD_ATTRIBUTE_DESKTOP_APP_POLICY</b> for the <i>Attribute</i> parameter:</p>
<p><b>PROCESS_CREATION_DESKTOP_APP_BREAKAWAY_ENABLE_PROCESS_TREE</b>                                         0x01</p>
<p>The process being created will create any child processes outside of the desktop app runtime environment.  This behavior is the default for processes for which no policy has been set.</p>
<p><b>PROCESS_CREATION_DESKTOP_APP_BREAKAWAY_DISABLE_PROCESS_TREE</b>                                           0x02</p>
<p>The process being created will create any child processes inside of the desktop app runtime environment.  This policy is inherited by the descendant processes until it is overridden by creating a process with <b>PROCESS_CREATION_DESKTOP_APP_BREAKAWAY_ENABLE_PROCESS_TREE</b>.</p>
<p><b>PROCESS_CREATION_DESKTOP_APP_BREAKAWAY_OVERRIDE</b>                                           0x04</p>
<p>The process being created will run inside the desktop app runtime environment.  This policy applies only to the process being created, not its descendants..</p>
<p>In order to launch the child process with the same protection level as the parent, the parent process must specify the <b>PROC_THREAD_ATTRIBUTE_PROTECTION_LEVEL</b> attribute for the child process. This can be used for both protected and unprotected processes. For example, when this flag is used by an unprotected process, the system will launch a child process at unprotected level. The <b>CREATE_PROTECTED_PROCESS</b> flag must be specified in both cases.</p>
<p>The following example launches a child process with the same protection level as the parent process:</p>
<pre><code class="lang-cpp">DWORD ProtectionLevel = PROTECTION_LEVEL_SAME;
SIZE_T AttributeListSize;

STARTUPINFOEXW StartupInfoEx = { 0 };

StartupInfoEx.StartupInfo.cb = sizeof(StartupInfoEx);

InitializeProcThreadAttributeList(NULL, 1, 0, &amp;AttributeListSize)


StartupInfoEx.lpAttributeList = (LPPROC_THREAD_ATTRIBUTE_LIST) HeapAlloc(
   GetProcessHeap(),
   0,
   AttributeListSize
   );

if (InitializeProcThreadAttributeList(StartupInfoEx.lpAttributeList,
                                     1,
                                     0,
                                     &amp;AttributeListSize) == FALSE)
{
   Result = GetLastError();
   goto exitFunc;
}

if (UpdateProcThreadAttribute(StartupInfoEx.lpAttributeList,
                             0,
                             PROC_THREAD_ATTRIBUTE_PROTECTION_LEVEL,
                             &amp;ProtectionLevel,
                             sizeof(ProtectionLevel),
                             NULL,
                             NULL) == FALSE)
{
   Result = GetLastError();
   goto exitFunc;
}

PROCESS_INFORMATION ProcessInformation = { 0 };

if (CreateProcessW(ApplicationName,
                  CommandLine,
                  ProcessAttributes,
                  ThreadAttributes,
                  InheritHandles,
                  EXTENDED_STARTUPINFO_PRESENT | CREATE_PROTECTED_PROCESS,
                  Environment,
                  CurrentDirectory,
                  (LPSTARTUPINFOW)&amp;StartupInfoEx,
                  &amp;ProcessInformation) == FALSE)
{
   Result = GetLastError();
   goto exitFunc;
}
</code></pre>
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
<td style="text-align: left;">Windows Vista [desktop apps only]</td>
</tr>
<tr>
<td><strong>Minimum supported server</strong></td>
<td style="text-align: left;">Windows Server 2008 [desktop apps only]</td>
</tr>
<tr>
<td><strong>Target Platform</strong></td>
<td style="text-align: left;">Windows</td>
</tr>
<tr>
<td><strong>Header</strong></td>
<td style="text-align: left;">processthreadsapi.h (include Windows.h on Windows 7, Windows Server 2008  Windows Server 2008 R2)</td>
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
<p><a href="/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-deleteprocthreadattributelist" data-linktype="absolute-path">DeleteProcThreadAttributeList</a></p>
<p><a href="/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-initializeprocthreadattributelist" data-linktype="absolute-path">InitializeProcThreadAttributeList</a></p>
<p><a href="/en-us/windows/desktop/ProcThread/process-and-thread-functions" data-linktype="absolute-path">Process and Thread Functions</a></p>
</dl></div>`

func TestGetAllSection(t *testing.T) {
	bufFile := bufio.NewReader(strings.NewReader(htm))
	main := symbols.GetMainContent(bufFile)
	fmt.Println(main.Length())
	sections := symbols.GetAllSection(main)
	fmt.Println("Done. Now printing")

	for k, li := range sections {
		fmt.Println(k)
		for _, l := range li {
			htm, er := goquery.OuterHtml(l)
			if er != nil {
				log.Panicln("Failed html conv")
			}
			fmt.Println(htm)
		}
		fmt.Println()
	}
	// desc, found := sections["basic-description"]
	// if found {
	// 	for _, l := range desc {
	// 		fmt.Println(goquery.OuterHtml(l))
	// 	}
	// }

}
