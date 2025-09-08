<script>
  import { onMount } from 'svelte';
  import { GetDeviceInventory } from '../../wailsjs/wailsjs/go/main/App.js';
  
  let devices = [];
  let isLoading = true;
  let sortBy = 'addedAt';
  let sortOrder = 'desc';
  let searchQuery = '';
  
  onMount(async () => {
    console.log('DeviceInventory mounted');
    await loadDevices();
  });
  
  async function loadDevices() {
    try {
      console.log('Loading device inventory from database...');
      isLoading = true;
      devices = await GetDeviceInventory();
      console.log('Device inventory loaded:', devices.length, 'devices');
    } catch (error) {
      console.error('Failed to load device inventory:', error);
      devices = [];
    } finally {
      isLoading = false;
    }
  }
  
  function sortDevices(field) {
    if (sortBy === field) {
      sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
    } else {
      sortBy = field;
      sortOrder = 'asc';
    }
  }
  
  function getSortIcon(field) {
    if (sortBy !== field) return '↕️';
    return sortOrder === 'asc' ? '↑' : '↓';
  }
  
  function formatDate(dateString) {
    if (!dateString) return 'Never';
    const date = new Date(dateString);
    return date.toLocaleString();
  }
  
  function getStatusColor(status) {
    switch (status) {
      case 'connected': return '#4CAF50';
      case 'disconnected': return '#9E9E9E';
      case 'testing': return '#FF9800';
      case 'failed': return '#F44336';
      default: return '#9E9E9E';
    }
  }
  
  function getTypeColor(type) {
    switch (type) {
      case 'eapi': return '#2196F3';
      case 'cloudvision': return '#9C27B0';
      case 'eos_rest': return '#FF9800';
      case 'telemetry': return '#4CAF50';
      default: return '#9E9E9E';
    }
  }
  
  $: filteredDevices = devices
    .filter(device => 
      device.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      device.url.toLowerCase().includes(searchQuery.toLowerCase()) ||
      device.type.toLowerCase().includes(searchQuery.toLowerCase()) ||
      device.status.toLowerCase().includes(searchQuery.toLowerCase())
    )
    .sort((a, b) => {
      let aVal = a[sortBy];
      let bVal = b[sortBy];
      
      // Handle date fields
      if (sortBy === 'addedAt' || sortBy === 'lastTested') {
        aVal = new Date(aVal).getTime();
        bVal = new Date(bVal).getTime();
      }
      
      // Handle string fields
      if (typeof aVal === 'string') {
        aVal = aVal.toLowerCase();
        bVal = bVal.toLowerCase();
      }
      
      if (sortOrder === 'asc') {
        return aVal > bVal ? 1 : -1;
      } else {
        return aVal < bVal ? 1 : -1;
      }
    });
</script>

<div class="device-inventory">
  <div class="section-header">
    <h2>Device Inventory</h2>
    <div class="header-actions">
      <input 
        type="text" 
        bind:value={searchQuery} 
        placeholder="Search devices..." 
        class="search-input"
      />
      <button class="refresh-button" on:click={loadDevices} disabled={isLoading}>
        {isLoading ? 'Loading...' : 'Refresh'}
      </button>
    </div>
  </div>

  {#if isLoading}
    <div class="loading">Loading device inventory...</div>
  {:else if filteredDevices.length === 0}
    <div class="no-devices">
      <p>No devices found in inventory.</p>
      <p>Add devices through the Endpoints tab to see them here.</p>
    </div>
  {:else}
    <div class="inventory-table">
      <table>
        <thead>
          <tr>
            <th class="sortable" on:click={() => sortDevices('name')}>
              Name {getSortIcon('name')}
            </th>
            <th class="sortable" on:click={() => sortDevices('type')}>
              Type {getSortIcon('type')}
            </th>
            <th class="sortable" on:click={() => sortDevices('url')}>
              URL {getSortIcon('url')}
            </th>
            <th class="sortable" on:click={() => sortDevices('username')}>
              Username {getSortIcon('username')}
            </th>
            <th class="sortable" on:click={() => sortDevices('status')}>
              Status {getSortIcon('status')}
            </th>
            <th class="sortable" on:click={() => sortDevices('addedAt')}>
              Added {getSortIcon('addedAt')}
            </th>
            <th class="sortable" on:click={() => sortDevices('lastTested')}>
              Last Tested {getSortIcon('lastTested')}
            </th>
            <th class="sortable" on:click={() => sortDevices('testCount')}>
              Tests {getSortIcon('testCount')}
            </th>
            <th class="sortable" on:click={() => sortDevices('successCount')}>
              Success Rate {getSortIcon('successCount')}
            </th>
            <th>Notes</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredDevices as device (device.id)}
            <tr>
              <td class="device-name">
                <strong>{device.name}</strong>
              </td>
              <td>
                <span class="type-badge" style="background-color: {getTypeColor(device.type)};">
                  {device.type.toUpperCase()}
                </span>
              </td>
              <td class="url-cell">
                <span class="url-text" title={device.url}>{device.url}</span>
              </td>
              <td>{device.username || '-'}</td>
              <td>
                <span class="status-badge" style="background-color: {getStatusColor(device.status)};">
                  {device.status}
                </span>
              </td>
              <td class="date-cell">{formatDate(device.addedAt)}</td>
              <td class="date-cell">{formatDate(device.lastTested)}</td>
              <td class="number-cell">{device.testCount}</td>
              <td class="number-cell">
                {device.testCount > 0 ? Math.round((device.successCount / device.testCount) * 100) : 0}%
              </td>
              <td class="notes-cell">
                <span class="notes-text" title={device.notes}>{device.notes || '-'}</span>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
    
    <div class="inventory-summary">
      <div class="summary-stats">
        <div class="stat">
          <span class="stat-label">Total Devices:</span>
          <span class="stat-value">{devices.length}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Connected:</span>
          <span class="stat-value">{devices.filter(d => d.status === 'connected').length}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Disconnected:</span>
          <span class="stat-value">{devices.filter(d => d.status === 'disconnected').length}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Failed:</span>
          <span class="stat-value">{devices.filter(d => d.status === 'failed').length}</span>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .device-inventory {
    padding: 1rem;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid rgba(255, 0, 128, 0.2);
  }

  .section-header h2 {
    margin: 0;
    color: #ff0080;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .header-actions {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .search-input {
    padding: 0.5rem;
    border: 1px solid rgba(255, 0, 128, 0.3);
    border-radius: 4px;
    background: rgba(0, 0, 0, 0.3);
    color: white;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .search-input:focus {
    outline: none;
    border-color: #ff0080;
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .refresh-button {
    padding: 0.5rem 1rem;
    background: linear-gradient(45deg, #ff0080, #ff4080);
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-weight: bold;
  }

  .refresh-button:hover:not(:disabled) {
    background: linear-gradient(45deg, #ff4080, #ff0080);
    box-shadow: 0 0 10px rgba(255, 0, 128, 0.3);
  }

  .refresh-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .loading, .no-devices {
    text-align: center;
    padding: 2rem;
    color: #e0e0e0;
  }

  .inventory-table {
    overflow-x: auto;
    margin-bottom: 1rem;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    background: rgba(0, 0, 0, 0.2);
    border-radius: 8px;
    overflow: hidden;
  }

  th, td {
    padding: 0.75rem;
    text-align: left;
    border-bottom: 1px solid rgba(255, 0, 128, 0.1);
  }

  th {
    background: rgba(255, 0, 128, 0.1);
    color: #ff0080;
    font-weight: bold;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .sortable {
    cursor: pointer;
    user-select: none;
  }

  .sortable:hover {
    background: rgba(255, 0, 128, 0.2);
  }

  td {
    color: #e0e0e0;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .device-name {
    font-weight: bold;
    color: #ff0080;
  }

  .type-badge, .status-badge {
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    color: white;
    font-size: 0.8rem;
    font-weight: bold;
    text-transform: uppercase;
  }

  .url-cell {
    max-width: 200px;
  }

  .url-text {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .date-cell {
    font-size: 0.9rem;
    color: #b0b0b0;
  }

  .number-cell {
    text-align: center;
    font-weight: bold;
  }

  .notes-cell {
    max-width: 150px;
  }

  .notes-text {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 0.9rem;
    color: #b0b0b0;
  }

  .inventory-summary {
    margin-top: 1rem;
    padding: 1rem;
    background: rgba(255, 0, 128, 0.05);
    border: 1px solid rgba(255, 0, 128, 0.2);
    border-radius: 8px;
  }

  .summary-stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 1rem;
  }

  .stat {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .stat-label {
    color: #e0e0e0;
    font-family: 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }

  .stat-value {
    color: #ff0080;
    font-weight: bold;
    font-family: 'Hermit', 'UbuntuMono', 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }
</style>