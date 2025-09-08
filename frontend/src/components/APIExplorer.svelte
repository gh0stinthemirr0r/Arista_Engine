<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { GetNetVisorAPIs, GetNetVisorAPIsByService, SearchNetVisorAPIs } from '../../wailsjs/wailsjs/go/main/App.js';
  
  const dispatch = createEventDispatcher();
  
  export let selectedEndpoint = null;
  export let apiCatalog = null;
  
  let selectedAPI = null;
  let requestMethod = "GET";
  let requestPath = "";
  let requestBody = "";
  let headers = "";
  let isLoading = false;
  let response = null;
  let netvisorAPIs = [];
  let searchQuery = "";
  let selectedService = "";
  
  onMount(async () => {
    console.log('APIExplorer mounted, loading NetVisor APIs...');
    await loadNetVisorAPIs();
  });

  async function loadNetVisorAPIs() {
    try {
      console.log('Loading APIs from NetVisor database...');
      netvisorAPIs = await GetNetVisorAPIs();
      console.log('NetVisor APIs loaded:', netvisorAPIs.length, 'APIs');
    } catch (error) {
      console.error('Failed to load NetVisor APIs:', error);
      netvisorAPIs = [];
    }
  }

  async function searchAPIs() {
    if (!searchQuery.trim()) {
      await loadNetVisorAPIs();
      return;
    }
    
    try {
      console.log('Searching NetVisor APIs for:', searchQuery);
      netvisorAPIs = await SearchNetVisorAPIs(searchQuery);
      console.log('Search results:', netvisorAPIs.length, 'APIs');
    } catch (error) {
      console.error('Failed to search NetVisor APIs:', error);
      netvisorAPIs = [];
    }
  }

  async function filterByService(service) {
    selectedService = service;
    try {
      if (service) {
        console.log('Filtering NetVisor APIs by service:', service);
        netvisorAPIs = await GetNetVisorAPIsByService(service);
      } else {
        await loadNetVisorAPIs();
      }
      console.log('Filtered APIs:', netvisorAPIs.length, 'APIs');
    } catch (error) {
      console.error('Failed to filter NetVisor APIs:', error);
      netvisorAPIs = [];
    }
  }

  function selectAPI(api) {
    selectedAPI = api;
    requestMethod = api.method || "GET";
    requestPath = api.path || "";
  }
  
  async function executeRequest() {
    if (!selectedEndpoint || !requestPath) return;
    
    isLoading = true;
    
    try {
      const { RunAPIRequest } = await import('../../wailsjs/wailsjs/go/main/App.js');
      const request = {
        endpointId: selectedEndpoint.id,
        method: requestMethod,
        path: requestPath,
        body: requestBody,
        headers: headers
      };
      
      response = await RunAPIRequest(request);
      dispatch('response', response);
    } catch (error) {
      console.error('API request failed:', error);
      response = {
        status: 500,
        headers: {},
        data: {
          error: error.message,
          message: "API request failed"
        }
      };
      dispatch('response', response);
    } finally {
      isLoading = false;
    }
  }
  
  function getMethodColor(method) {
    switch (method?.toUpperCase()) {
      case 'GET': return '#4CAF50';
      case 'POST': return '#2196F3';
      case 'PUT': return '#FF9800';
      case 'DELETE': return '#F44336';
      case 'RUNCMDS': return '#9C27B0';
      default: return '#9E9E9E';
    }
  }
</script>

<div class="api-explorer">
  <div class="section-header">
    <h2>API Explorer</h2>
    {#if selectedEndpoint}
      <div class="endpoint-info">
        <span class="endpoint-name">{selectedEndpoint.name}</span>
        <span class="endpoint-type">{selectedEndpoint.type}</span>
      </div>
    {/if}
  </div>
  
  {#if !selectedEndpoint}
    <div class="no-endpoint">
      <p>Please select an endpoint first</p>
    </div>
  {:else}
    <div class="explorer-content">
      <div class="api-selection">
        <div class="api-controls">
          <h3>NetVisor API Database</h3>
          <div class="search-controls">
            <input 
              type="text" 
              bind:value={searchQuery} 
              placeholder="Search APIs..." 
              class="search-input"
              on:input={searchAPIs}
            />
            <select bind:value={selectedService} on:change={() => filterByService(selectedService)} class="service-filter">
              <option value="">All Services</option>
              <option value="eapi">EOS eAPI</option>
              <option value="cloudvision">CloudVision</option>
              <option value="eos_rest">EOS REST</option>
              <option value="telemetry">Telemetry</option>
            </select>
            <button class="refresh-btn" on:click={loadNetVisorAPIs}>Refresh</button>
          </div>
        </div>
        
        <div class="api-categories">
          {#if netvisorAPIs.length > 0}
            <div class="category-section">
              <h4>Available APIs ({netvisorAPIs.length})</h4>
              <div class="api-list">
                {#each netvisorAPIs as api}
                  <div class="api-item" on:click={() => selectAPI(api)}>
                    <span class="method" style="color: {getMethodColor(api.method)}">{api.method || 'GET'}</span>
                    <span class="path">{api.path || 'N/A'}</span>
                    <span class="description">{api.description || 'No description available'}</span>
                    {#if api.category}
                      <span class="category">{api.category}</span>
                    {/if}
                    {#if api.tags}
                      <span class="tags">{api.tags}</span>
                    {/if}
                  </div>
                {/each}
              </div>
            </div>
          {:else}
            <div class="no-apis">
              <p>No APIs found in NetVisor database</p>
              <p>Make sure the netvisor_api_v711.db file exists in the project root</p>
            </div>
          {/if}
        </div>
      </div>
      
      <div class="request-builder">
        <h3>Request Builder</h3>
        <div class="request-form">
          <div class="form-row">
            <label>Method</label>
            <select bind:value={requestMethod}>
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
              <option value="runCmds">runCmds (eAPI)</option>
            </select>
          </div>
          
          <div class="form-row">
            <label>Path</label>
            <input type="text" bind:value={requestPath} placeholder="/api/v1/..." />
          </div>
          
          {#if requestMethod === 'POST' || requestMethod === 'PUT'}
            <div class="form-row">
              <label>Request Body</label>
              <textarea bind:value={requestBody} placeholder="JSON request body..."></textarea>
            </div>
          {/if}
          
          <div class="form-row">
            <label>Headers (JSON)</label>
              <textarea bind:value={headers} placeholder="Content-Type: application/json"></textarea>
          </div>
          
          <div class="form-actions">
            <button class="execute-btn" on:click={executeRequest} disabled={isLoading || !requestPath}>
              {isLoading ? 'Executing...' : 'Execute Request'}
            </button>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .api-explorer {
    max-width: 1400px;
    margin: 0 auto;
    padding: 1rem;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid rgba(255, 0, 128, 0.2);
  }

  .section-header h2 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    color: #ff0080;
    margin: 0;
  }

  .endpoint-info {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .endpoint-name {
    color: #ff0080;
    font-weight: bold;
  }

  .endpoint-type {
    background: rgba(255, 0, 128, 0.2);
    color: #ff0080;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.8rem;
    text-transform: uppercase;
  }

  .no-endpoint {
    text-align: center;
    padding: 3rem;
    color: #888;
  }

  .explorer-content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }

  .api-controls {
    margin-bottom: 1.5rem;
  }

  .api-controls h3 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    color: #ff0080;
    margin: 0 0 1rem 0;
  }

  .search-controls {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .search-input, .service-filter {
    padding: 0.5rem;
    border: 1px solid rgba(255, 0, 128, 0.3);
    border-radius: 4px;
    background: rgba(0, 0, 0, 0.3);
    color: white;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .search-input:focus, .service-filter:focus {
    outline: none;
    border-color: #ff0080;
    box-shadow: 0 0 5px rgba(255, 0, 128, 0.3);
  }

  .refresh-btn {
    background: linear-gradient(45deg, #ff0080, #ff4080);
    border: 1px solid #ff0080;
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
  }

  .refresh-btn:hover {
    background: linear-gradient(45deg, #ff4080, #ff0080);
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .api-categories {
    max-height: 600px;
    overflow-y: auto;
  }

  .category-section h4 {
    color: #ff0080;
    margin: 0 0 1rem 0;
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
  }

  .api-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .api-item {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    padding: 1rem;
    cursor: pointer;
    transition: all 0.3s ease;
    display: grid;
    grid-template-columns: auto 1fr auto;
    gap: 1rem;
    align-items: center;
  }

  .api-item:hover {
    border-color: rgba(255, 0, 128, 0.3);
    background: rgba(255, 255, 255, 0.08);
  }

  .method {
    font-weight: bold;
    font-size: 0.8rem;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    background: rgba(0, 0, 0, 0.3);
    min-width: 60px;
    text-align: center;
  }

  .path {
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    color: #e0e0e0;
    font-size: 0.9rem;
  }

  .description {
    color: #ccc;
    font-size: 0.85rem;
    grid-column: 1 / -1;
    margin-top: 0.5rem;
  }

  .category, .tags {
    font-size: 0.7rem;
    color: #888;
    background: rgba(255, 255, 255, 0.1);
    padding: 0.2rem 0.4rem;
    border-radius: 3px;
    margin-left: 0.5rem;
  }

  .no-apis {
    text-align: center;
    padding: 2rem;
    color: #888;
  }

  .request-builder {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 0, 128, 0.2);
    border-radius: 8px;
    padding: 1.5rem;
  }

  .request-builder h3 {
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    color: #ff0080;
    margin: 0 0 1.5rem 0;
  }

  .request-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .form-row {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-row label {
    color: #ccc;
    font-size: 0.9rem;
    font-weight: bold;
  }

  .form-row input,
  .form-row select,
  .form-row textarea {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: white;
    padding: 0.75rem;
    border-radius: 4px;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .form-row input:focus,
  .form-row select:focus,
  .form-row textarea:focus {
    outline: none;
    border-color: #ff0080;
    box-shadow: 0 0 5px rgba(255, 0, 128, 0.3);
  }

  .form-row textarea {
    min-height: 100px;
    resize: vertical;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
  }

  .execute-btn {
    background: linear-gradient(45deg, #ff0080, #ff4080);
    border: 1px solid #ff0080;
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    font-family: 'Hermit', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
  }

  .execute-btn:hover:not(:disabled) {
    background: linear-gradient(45deg, #ff4080, #ff0080);
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .execute-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>