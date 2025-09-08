<script>
  export let response = null;
  
  let viewMode = 'json';
  
  function formatJSON(obj) {
    return JSON.stringify(obj, null, 2);
  }
</script>

<div class="response-viewer">
  <div class="section-header">
    <h2>Response Viewer</h2>
    {#if response}
      <div class="response-info">
        <span class="status-code" class:success={response.status >= 200 && response.status < 300}>
          {response.status}
        </span>
        <span class="response-time">~1000ms</span>
      </div>
    {/if}
  </div>
  
  {#if !response}
    <div class="no-response">
      <p>No response data available</p>
    </div>
  {:else}
    <div class="response-content">
      <div class="view-controls">
        <button 
          class="view-button" 
          class:active={viewMode === 'json'}
          on:click={() => viewMode = 'json'}
        >
          JSON
        </button>
        <button 
          class="view-button" 
          class:active={viewMode === 'raw'}
          on:click={() => viewMode = 'raw'}
        >
          Raw
        </button>
        <button 
          class="view-button" 
          class:active={viewMode === 'headers'}
          on:click={() => viewMode = 'headers'}
        >
          Headers
        </button>
      </div>
      
      <div class="response-data">
        {#if viewMode === 'json'}
          <pre class="json-viewer">{formatJSON(response.data)}</pre>
        {:else if viewMode === 'raw'}
          <pre class="raw-viewer">{JSON.stringify(response, null, 2)}</pre>
        {:else if viewMode === 'headers'}
          <div class="headers-viewer">
            {#each Object.entries(response.headers) as [key, value]}
              <div class="header-row">
                <span class="header-key">{key}:</span>
                <span class="header-value">{value}</span>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .response-viewer {
    max-width: 1200px;
    margin: 0 auto;
  }
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }
  
  .section-header h2 {
    margin: 0;
    color: #ff0080;
  }
  
  .response-info {
    display: flex;
    gap: 1rem;
    align-items: center;
  }
  
  .status-code {
    background: rgba(255, 0, 128, 0.2);
    color: #ff0080;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-weight: bold;
  }
  
  .status-code.success {
    background: rgba(76, 175, 80, 0.2);
    color: #4CAF50;
  }
  
  .response-time {
    color: #888;
    font-size: 0.9rem;
  }
  
  .no-response {
    text-align: center;
    padding: 3rem;
    color: #888;
  }
  
  .response-content {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .view-controls {
    display: flex;
    gap: 0.5rem;
  }
  
  .view-button {
    background: transparent;
    border: 1px solid rgba(255, 0, 128, 0.3);
    color: #ccc;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-family: inherit;
    transition: all 0.3s ease;
  }
  
  .view-button:hover {
    border-color: rgba(255, 0, 128, 0.5);
    color: #ff0080;
  }
  
  .view-button.active {
    background: rgba(255, 0, 128, 0.2);
    border-color: #ff0080;
    color: #ff0080;
  }
  
  .response-data {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 0, 128, 0.2);
    border-radius: 8px;
    padding: 1rem;
    min-height: 400px;
    overflow: auto;
  }
  
  .json-viewer,
  .raw-viewer {
    margin: 0;
    color: #ccc;
    font-family: 'Courier New', monospace;
    font-size: 0.9rem;
    line-height: 1.4;
    white-space: pre-wrap;
    word-break: break-word;
  }
  
  .headers-viewer {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .header-row {
    display: flex;
    gap: 1rem;
    padding: 0.5rem;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
  }
  
  .header-key {
    color: #ff0080;
    font-weight: bold;
    min-width: 150px;
  }
  
  .header-value {
    color: #ccc;
    font-family: monospace;
  }
</style>