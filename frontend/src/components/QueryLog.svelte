<script>
  export let queryLog = [];
  
  let selectedLog = null;
  let filterStatus = 'all';
  let searchQuery = '';
  
  $: filteredLogs = getFilteredLogs();
  
  function getFilteredLogs() {
    let logs = queryLog;
    
    if (filterStatus !== 'all') {
      logs = logs.filter(log => log.status === filterStatus);
    }
    
    if (searchQuery) {
      logs = logs.filter(log => 
        log.endpoint?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        log.method?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        log.path?.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }
    
    return logs.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));
  }
  
  function getStatusColor(status) {
    switch (status) {
      case 'success': return '#4CAF50';
      case 'error': return '#F44336';
      case 'pending': return '#FF9800';
      default: return '#9E9E9E';
    }
  }
  
  function formatTimestamp(timestamp) {
    return new Date(timestamp).toLocaleString();
  }
  
  function clearLogs() {
    queryLog = [];
  }
</script>

<div class="query-log">
  <div class="section-header">
    <h2>Query Log</h2>
    <div class="log-actions">
      <span class="log-count">{filteredLogs.length} entries</span>
      <button class="clear-button" on:click={clearLogs}>Clear Logs</button>
    </div>
  </div>
  
  <div class="log-controls">
    <div class="search-box">
      <input 
        type="text" 
        bind:value={searchQuery} 
        placeholder="Search logs..." 
        class="search-input"
      />
    </div>
    <div class="status-filter">
      <select bind:value={filterStatus}>
        <option value="all">All Status</option>
        <option value="success">Success</option>
        <option value="error">Error</option>
        <option value="pending">Pending</option>
      </select>
    </div>
  </div>
  
  <div class="log-list">
    {#each filteredLogs as log}
      <div 
        class="log-card" 
        class:selected={selectedLog?.id === log.id}
        on:click={() => selectedLog = log}
      >
        <div class="log-header">
          <div class="log-method" style="color: {getStatusColor(log.status)}">
            {log.method || 'GET'}
          </div>
          <div class="log-path">{log.path || '/api/endpoint'}</div>
          <div class="log-status" style="color: {getStatusColor(log.status)}">
            {log.status || 'unknown'}
          </div>
        </div>
        <div class="log-details">
          <div class="log-endpoint">{log.endpoint || 'Unknown Endpoint'}</div>
          <div class="log-timestamp">{formatTimestamp(log.timestamp || new Date())}</div>
        </div>
      </div>
    {:else}
      <div class="no-logs">
        <p>No log entries found</p>
      </div>
    {/each}
  </div>
  
  <!-- Log Detail Modal -->
  {#if selectedLog}
    <div class="log-modal-overlay" on:click={() => selectedLog = null}>
      <div class="log-modal" on:click|stopPropagation>
        <div class="modal-header">
          <h3>Query Details</h3>
          <button class="close-button" on:click={() => selectedLog = null}>×</button>
        </div>
        <div class="modal-content">
          <div class="detail-section">
            <h4>Request</h4>
            <div class="detail-grid">
              <div class="detail-item">
                <label>Method:</label>
                <span>{selectedLog.method || 'GET'}</span>
              </div>
              <div class="detail-item">
                <label>Path:</label>
                <span>{selectedLog.path || '/api/endpoint'}</span>
              </div>
              <div class="detail-item">
                <label>Endpoint:</label>
                <span>{selectedLog.endpoint || 'Unknown'}</span>
              </div>
              <div class="detail-item">
                <label>Status:</label>
                <span style="color: {getStatusColor(selectedLog.status)}">{selectedLog.status || 'unknown'}</span>
              </div>
            </div>
          </div>
          
          {#if selectedLog.requestBody}
            <div class="detail-section">
              <h4>Request Body</h4>
              <pre class="code-block">{JSON.stringify(selectedLog.requestBody, null, 2)}</pre>
            </div>
          {/if}
          
          {#if selectedLog.response}
            <div class="detail-section">
              <h4>Response</h4>
              <pre class="code-block">{JSON.stringify(selectedLog.response, null, 2)}</pre>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .query-log {
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
  
  .log-actions {
    display: flex;
    gap: 1rem;
    align-items: center;
  }
  
  .log-count {
    color: #888;
    font-size: 0.9rem;
  }
  
  .clear-button {
    background: transparent;
    border: 1px solid #666;
    color: #ccc;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-family: inherit;
    transition: all 0.3s ease;
  }
  
  .clear-button:hover {
    border-color: #ff4444;
    color: #ff4444;
  }
  
  .log-controls {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
  }
  
  .search-box {
    flex: 1;
  }
  
  .search-input {
    width: 100%;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 0, 128, 0.3);
    border-radius: 4px;
    padding: 0.75rem;
    color: white;
    font-family: inherit;
  }
  
  .search-input:focus {
    outline: none;
    border-color: #ff0080;
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.2);
  }
  
  .status-filter select {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 0, 128, 0.3);
    border-radius: 4px;
    padding: 0.75rem;
    color: white;
    font-family: inherit;
    min-width: 150px;
  }
  
  .status-filter select:focus {
    outline: none;
    border-color: #ff0080;
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.2);
  }
  
  .log-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .log-card {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 0, 128, 0.2);
    border-radius: 8px;
    padding: 1rem;
    cursor: pointer;
    transition: all 0.3s ease;
  }
  
  .log-card:hover {
    border-color: rgba(255, 0, 128, 0.5);
    background: rgba(255, 0, 128, 0.1);
  }
  
  .log-card.selected {
    border-color: #ff0080;
    background: rgba(255, 0, 128, 0.2);
  }
  
  .log-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 0.5rem;
  }
  
  .log-method {
    font-weight: bold;
    font-size: 0.8rem;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    background: rgba(255, 255, 255, 0.1);
    min-width: 50px;
    text-align: center;
  }
  
  .log-path {
    color: #ccc;
    font-family: monospace;
    flex: 1;
  }
  
  .log-status {
    font-weight: bold;
    text-transform: uppercase;
    font-size: 0.8rem;
  }
  
  .log-details {
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: #888;
    font-size: 0.9rem;
  }
  
  .no-logs {
    text-align: center;
    padding: 3rem;
    color: #888;
  }
  
  .log-modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  
  .log-modal {
    background: #1a1a1a;
    border: 1px solid rgba(255, 0, 128, 0.3);
    border-radius: 8px;
    max-width: 800px;
    max-height: 80vh;
    width: 90%;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    border-bottom: 1px solid rgba(255, 0, 128, 0.3);
  }
  
  .modal-header h3 {
    margin: 0;
    color: #ff0080;
  }
  
  .close-button {
    background: none;
    border: none;
    color: #ccc;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 0;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .close-button:hover {
    color: #ff0080;
  }
  
  .modal-content {
    padding: 1rem;
    overflow: auto;
    flex: 1;
  }
  
  .detail-section {
    margin-bottom: 1.5rem;
  }
  
  .detail-section h4 {
    margin: 0 0 1rem 0;
    color: #ff0080;
  }
  
  .detail-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
  }
  
  .detail-item {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .detail-item label {
    color: #888;
    font-size: 0.9rem;
  }
  
  .detail-item span {
    color: #ccc;
    font-family: monospace;
  }
  
  .code-block {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 0, 128, 0.2);
    border-radius: 4px;
    padding: 1rem;
    color: #ccc;
    font-family: 'Courier New', monospace;
    font-size: 0.9rem;
    line-height: 1.4;
    white-space: pre-wrap;
    word-break: break-word;
    max-height: 300px;
    overflow: auto;
  }
</style>