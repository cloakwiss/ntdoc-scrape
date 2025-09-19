package function_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/cloakwiss/ntdocs/utils"
	"github.com/k0kubun/pp/v3"
)

func TestParameters(t *testing.T) {
	var data string = `<div class="content"><p>Releases, decommits, or releases and decommits a region of pages within the virtual address space of the calling process.</p>
<p>To free memory allocated in another process by the <a href="/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualallocex" data-linktype="absolute-path">VirtualAllocEx</a> function, use the <a href="/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualfreeex" data-linktype="absolute-path">VirtualFreeEx</a> function.</p>
<h2 id="syntax">Syntax</h2>
<pre><code class="lang-cpp">BOOL VirtualFree(
  [in] LPVOID lpAddress,
  [in] SIZE_T dwSize,
  [in] DWORD  dwFreeType
);
</code></pre>
<h2 id="parameters">Parameters</h2>
<p><code>[in] lpAddress</code></p>
<p>A pointer to the base address of the region of pages to be freed.</p>
<p>If the <em>dwFreeType</em> parameter is <strong>MEM_RELEASE</strong>, this parameter must be the base address returned by the <a href="/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc" data-linktype="absolute-path">VirtualAlloc</a> function when the region of pages is reserved.</p>
<p><code>[in] dwSize</code></p>
<p>The size of the region of memory to be freed, in bytes.</p>
<p>If the <em>dwFreeType</em> parameter is <strong>MEM_RELEASE</strong>, this parameter must be 0 (zero). The function frees the entire region that is reserved in the initial allocation call to <a href="/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc" data-linktype="absolute-path">VirtualAlloc</a>.</p>
<p>If the <em>dwFreeType</em> parameter is <strong>MEM_DECOMMIT</strong>, the function decommits all memory pages that contain one or more bytes in the range from the <em>lpAddress</em> parameter to <code>(lpAddress+dwSize)</code>. This means, for example, that a 2-byte region of memory that straddles a page boundary causes both pages to be decommitted. If <em>lpAddress</em> is the base address returned by <a href="/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc" data-linktype="absolute-path">VirtualAlloc</a> and <em>dwSize</em> is 0 (zero), the function decommits the entire region that is allocated by <strong>VirtualAlloc</strong>. After that, the entire region is in the reserved state.</p>
<p><code>[in] dwFreeType</code></p>
<p>The type of free operation. This parameter must be one of the following values.</p>
<table>
<tbody><tr>
<th>Value</th>
<th>Meaning</th>
</tr>
<tr>
<td width="40%"><a id="MEM_DECOMMIT"></a><a id="mem_decommit"></a><dl>
<dt><b>MEM_DECOMMIT</b></dt>
<dt>0x00004000</dt>
</dl>
</td>
<td width="60%">
Decommits the specified region of committed pages. After the operation, the pages are in the reserved state.
<p>The function does not fail if you attempt to decommit an uncommitted page. This means that you can decommit a range of pages without first determining the current commitment state.</p>
<p>The <strong>MEM_DECOMMIT</strong> value is not supported when the <em>lpAddress</em> parameter provides the base address for an enclave. This is true for enclaves that do not support dynamic memory management (i.e. SGX1).  SGX2 enclaves permit <strong>MEM_DECOMMIT</strong> anywhere in the enclave.</p>
</td>
</tr>
<tr>
<td width="40%"><a id="MEM_RELEASE"></a><a id="mem_release"></a><dl>
<dt><b>MEM_RELEASE</b></dt>
<dt>0x00008000</dt>
</dl>
</td>
<td width="60%">
Releases the specified region of pages, or placeholder (for a placeholder, the address space is released and available for other allocations). After this operation, the pages are in the free state.
<p>If you specify this value, <em>dwSize</em> must be 0 (zero), and <em>lpAddress</em> must point to the base address returned by the <a href="/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc" data-linktype="absolute-path">VirtualAlloc</a> function when the region is reserved. The function fails if either of these conditions is not met.</p>
<p>If any pages in the region are committed currently, the function first decommits, and then releases them.</p>
<p>The function does not fail if you attempt to release pages that are in different states, some reserved and some committed. This means that you can release a range of pages without first determining the current commitment state.</p>
</td>
</tr>
</tbody></table>
<p>When using <strong>MEM_RELEASE</strong>, this parameter can additionally specify one of the following values.</p>
<table>
<tbody><tr>
<th>Value</th>
<th>Meaning</th>
</tr>
<tr>
<td width="40%"><a id="MEM_COALESCE_PLACEHOLDERS"></a><a id="mem_coalesce_placeholders"></a><dl>
<dt><b>MEM_COALESCE_PLACEHOLDERS</b></dt>
<dt>0x00000001</dt>
</dl>
</td>
<td width="60%">
To coalesce two adjacent placeholders, specify <code>MEM_RELEASE | MEM_COALESCE_PLACEHOLDERS</code>. When you coalesce placeholders, <i>lpAddress</i> and <i>dwSize</i> must exactly match the overall range of the placeholders to be merged.
</td>
</tr>
<tr>
<td width="40%"><a id="MEM_PRESERVE_PLACEHOLDER"></a><a id="mem_preserve_placeholder"></a><dl>
<dt><b>MEM_PRESERVE_PLACEHOLDER</b></dt>
<dt>0x00000002</dt>
</dl>
</td>
<td width="60%">
Frees an allocation back to a placeholder (after you've replaced a placeholder with a private allocation using <a href="nf-memoryapi-virtualalloc2" data-linktype="relative-path">VirtualAlloc2</a> or <a href="https://msdn.microsoft.com/en-us/library/Mt832850(v=VS.85).aspx" data-linktype="external">Virtual2AllocFromApp</a>).
<p>To split a placeholder into two placeholders, specify <code>MEM_RELEASE | MEM_PRESERVE_PLACEHOLDER</code>.</p>
</td>
</tr>
</tbody></table>
<h2 id="return-value">Return value</h2>
<p>If the function succeeds, the return value is nonzero.</p>
<p>If the function fails, the return value is 0 (zero). To get extended error information, call <a href="/en-us/windows/win32/api/errhandlingapi/nf-errhandlingapi-getlasterror" data-linktype="absolute-path">GetLastError</a>.</p>
<h2 id="remarks">Remarks</h2>
<p>Each page of memory in a process virtual address space has a <a href="/en-us/windows/win32/Memory/page-state" data-linktype="absolute-path">Page State</a>. The <strong>VirtualFree</strong> function can decommit a range of pages that are in different states, some committed and some uncommitted. This means that you can decommit a range of pages without first determining the current commitment state of each page. Decommitting a page releases its physical storage, either in memory or in the paging file on disk.</p>
<p>If a page is decommitted but not released, its state changes to reserved. Subsequently, you can call <a href="/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc" data-linktype="absolute-path">VirtualAlloc</a> to commit it, or <strong>VirtualFree</strong> to release it. Attempts to read from or write to a reserved page results in an access violation exception.</p>
<p>The <strong>VirtualFree</strong> function can release a range of pages that are in different states, some reserved and some committed. This means that you can release a range of pages without first determining the current commitment state of each page. The entire range of pages originally reserved by the <a href="nf-memoryapi-virtualalloc" data-linktype="relative-path">VirtualAlloc</a> function must be released at the same time.</p>
<p>If a page is released, its state changes to free, and it is available for subsequent allocation operations. After memory is released or decommited, you can never refer to the memory again. Any information that may have been in that memory is gone forever. Attempting to read from or write to a free page results in an access violation exception. If you need to keep information, do not decommit or free memory that contains the information.</p>
<p>The <strong>VirtualFree</strong> function can be used on an AWE region of memory, and it invalidates any physical page mappings in the region when freeing the address space. However, the physical page is not deleted, and the application can use them. The application must explicitly call <a href="nf-memoryapi-freeuserphysicalpages" data-linktype="relative-path">FreeUserPhysicalPages</a> to free the physical pages. When the process is terminated, all resources are cleaned up automatically.</p>
<p><strong>Windows&nbsp;10, version 1709 and later and Windows 11:</strong> To delete the enclave when you finish using it, call <a href="../enclaveapi/nf-enclaveapi-deleteenclave" data-linktype="relative-path">DeleteEnclave</a>. You cannot delete a VBS enclave by calling the <strong>VirtualFree</strong> or <a href="nf-memoryapi-virtualfreeex" data-linktype="relative-path">VirtualFreeEx</a> function. You can still delete an SGX enclave by calling <strong>VirtualFree</strong> or <strong>VirtualFreeEx</strong>.</p>
<p><strong>Windows&nbsp;10, version 1507, Windows&nbsp;10, version 1511, Windows&nbsp;10, version 1607 and Windows&nbsp;10, version 1703:</strong> To delete the enclave when you finish using it, call the <strong>VirtualFree</strong> or <a href="nf-memoryapi-virtualfreeex" data-linktype="relative-path">VirtualFreeEx</a> function and specify the following values:</p>
<ul>
<li>The base address of the enclave for the <em>lpAddress</em> parameter.</li>
<li>0 for the <em>dwSize</em> parameter.</li>
<li><strong>MEM_RELEASE</strong> for the <em>dwFreeType</em> parameter.</li>
</ul>
<h3 id="examples">Examples</h3>
<p>For an example, see <a href="/en-us/windows/win32/Memory/reserving-and-committing-memory" data-linktype="absolute-path">Reserving and Committing Memory</a>.</p>
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
<td style="text-align: left;">Windows&nbsp;XP [desktop apps | UWP apps]</td>
</tr>
<tr>
<td><strong>Minimum supported server</strong></td>
<td style="text-align: left;">Windows Server&nbsp;2003 [desktop apps | UWP apps]</td>
</tr>
<tr>
<td><strong>Target Platform</strong></td>
<td style="text-align: left;">Windows</td>
</tr>
<tr>
<td><strong>Header</strong></td>
<td style="text-align: left;">memoryapi.h (include Windows.h, Memoryapi.h)</td>
</tr>
<tr>
<td><strong>Library</strong></td>
<td style="text-align: left;">onecore.lib</td>
</tr>
<tr>
<td><strong>DLL</strong></td>
<td style="text-align: left;">Kernel32.dll</td>
</tr>
</tbody>
</table>
<h2 id="see-also">See also</h2>
<p><a href="/en-us/windows/win32/Memory/memory-management-functions" data-linktype="absolute-path">Memory Management Functions</a></p>
<p><a href="/en-us/windows/win32/Memory/virtual-memory-functions" data-linktype="absolute-path">Virtual Memory Functions</a></p>
<p><a href="nf-memoryapi-virtualfreeex" data-linktype="relative-path">VirtualFreeEx</a></p>
<p><a href="/en-us/windows/win32/trusted-execution/enclaves-available-in-vertdll" data-linktype="absolute-path">Vertdll APIs available in VBS enclaves</a></p>
</div>`

	backing := strings.NewReader(data)
	buffer := bufio.NewReader(backing)
	mainContent := utils.GetMainContent(buffer)
	content := utils.GetAllSection(mainContent)
	// paras, er := function.HandleParameterSectionOfFunction(content["parameters"])
	// if er == nil {
	// 	pp.Println(paras)
	// }

	requirements, er := utils.HandleRequriementSectionOfFunction(content["requirements"])
	if er == nil {
		pp.Println(requirements)
	}
}
