// Configuration
const API_BASE_URL = window.location.protocol === "https:" ? "https://api.hookwatch.antcoders.dev" : "http://api.hookwatch.antcoders.dev";

// Utility functions
function showResult(elementId, content, type = "info") {
  const element = document.getElementById(elementId);
  element.innerHTML = content;
  element.className = `result ${type}`;
}

function showLoading(elementId) {
  const element = document.getElementById(elementId);
  element.innerHTML = '<div class="loading">Loading...</div>';
  element.className = "result";
}

function formatJSON(obj) {
  return JSON.stringify(obj, null, 2);
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleString();
}

// API functions
async function makeRequest(url, options = {}) {
  try {
    const response = await fetch(url, {
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || `HTTP ${response.status}: ${response.statusText}`);
    }

    return data;
  } catch (error) {
    throw new Error(`Request failed: ${error.message}`);
  }
}

// Send webhook
async function sendWebhook() {
  const endpointId = document.getElementById("endpointId").value.trim();
  const httpMethod = document.getElementById("httpMethod").value;
  const webhookData = document.getElementById("webhookData").value.trim();

  if (!endpointId) {
    showResult("sendResult", "Please enter an endpoint ID", "error");
    return;
  }

  if (httpMethod !== "GET" && !webhookData) {
    showResult("sendResult", "Please enter webhook data", "error");
    return;
  }

  let parsedData;
  let requestOptions = {
    method: httpMethod,
    headers: {
      "Content-Type": "application/json",
    },
  };

  // Handle different HTTP methods
  if (httpMethod === "GET") {
    // For GET requests, use query parameters instead of body
    if (webhookData) {
      try {
        parsedData = JSON.parse(webhookData);
        const params = new URLSearchParams();
        for (const [key, value] of Object.entries(parsedData)) {
          if (typeof value === "object") {
            params.append(key, JSON.stringify(value));
          } else {
            params.append(key, value);
          }
        }
        const queryString = params.toString();
        if (queryString) {
          requestOptions.url = `${API_BASE_URL}/webhooks/${endpointId}/receive?${queryString}`;
        } else {
          requestOptions.url = `${API_BASE_URL}/webhooks/${endpointId}/receive`;
        }
      } catch (error) {
        showResult("sendResult", `Invalid JSON for query parameters: ${error.message}`, "error");
        return;
      }
    } else {
      requestOptions.url = `${API_BASE_URL}/webhooks/${endpointId}/receive`;
    }
    // Remove Content-Type for GET requests
    delete requestOptions.headers["Content-Type"];
  } else {
    // For other methods, use JSON body
    if (webhookData) {
      try {
        parsedData = JSON.parse(webhookData);
        requestOptions.body = JSON.stringify(parsedData);
      } catch (error) {
        showResult("sendResult", `Invalid JSON: ${error.message}`, "error");
        return;
      }
    } else {
      requestOptions.body = JSON.stringify({});
    }
    requestOptions.url = `${API_BASE_URL}/webhooks/${endpointId}/receive`;
  }

  showLoading("sendResult");

  try {
    const result = await makeRequest(requestOptions.url, requestOptions);

    const successMessage = `
      <h4>✅ Webhook sent successfully!</h4>
      <div class="meta">
        <strong>Endpoint:</strong> ${result.endpointId}<br>
        <strong>Method:</strong> ${result.method}<br>
        <strong>Webhook ID:</strong> ${result.webhookId}<br>
        <strong>Status:</strong> ${result.status}<br>
        <strong>Timestamp:</strong> ${formatDate(result.timestamp)}
      </div>
      <div class="data">
        <strong>Response:</strong><br>
        ${formatJSON(result)}
      </div>
    `;

    showResult("sendResult", successMessage, "success");

    // Auto-load logs after sending
    setTimeout(() => {
      document.getElementById("logsEndpointId").value = endpointId;
      loadLogs();
    }, 1000);
  } catch (error) {
    showResult("sendResult", `❌ Error: ${error.message}`, "error");
  }
}

// Load webhook logs
async function loadLogs() {
  const endpointId = document.getElementById("logsEndpointId").value.trim();

  if (!endpointId) {
    showResult("logsResult", "Please enter an endpoint ID", "error");
    return;
  }

  showLoading("logsResult");

  try {
    const result = await makeRequest(`${API_BASE_URL}/webhooks/${endpointId}/logs?limit=20`);

    // Check if logs array exists and has items
    if (!result.logs || result.logs.length === 0) {
      showResult(
        "logsResult",
        `
        <div class="info">
          <h4>📭 No webhooks found</h4>
          <p>No webhook logs found for endpoint: <strong>${endpointId}</strong></p>
          <p>Try sending a webhook first using the form above!</p>
        </div>
      `,
        "info"
      );
      return;
    }

    let logsHtml = `
      <h4>📋 Webhook Logs (${result.count || result.logs.length} found)</h4>
      <div class="meta">
        <strong>Endpoint:</strong> ${result.endpointId} | 
        <strong>Limit:</strong> ${result.limit || 20}
      </div>
    `;

    result.logs.forEach((log, index) => {
      const statusClass = `status-${log.status}`;
      const statusBadge = `<span class="status-badge ${statusClass}">${log.status}</span>`;
      const logId = `log-${index}`;

      logsHtml += `
                <div class="webhook-log">
                    <h4>Webhook #${index + 1} - ${statusBadge}</h4>
                    <div class="meta">
                        <strong>ID:</strong> ${log.id} | 
                        <strong>Method:</strong> ${log.method} | 
                        <strong>IP:</strong> ${log.ip_address} | 
                        <strong>Created:</strong> ${formatDate(log.created_at)}
                        <button onclick="showReplayDialog('${log.id}', '${log.method}')" class="btn btn-info btn-sm">🔄 Replay</button>
                    </div>
                    <div class="data-sections">
                        <div class="data-section">
                            <div class="data-header collapsed" onclick="toggleSection('${logId}-headers')">
                                <span class="toggle-icon">▶</span>
                                <strong>Headers</strong>
                            </div>
                            <div class="data-content" id="${logId}-headers" style="display: none;">
                                ${formatJSON(log.headers)}
                            </div>
                        </div>
                        <div class="data-section">
                            <div class="data-header collapsed" onclick="toggleSection('${logId}-body')">
                                <span class="toggle-icon">▶</span>
                                <strong>Body</strong>
                            </div>
                            <div class="data-content" id="${logId}-body" style="display: none;">
                                ${formatJSON(log.body)}
                            </div>
                        </div>
                    </div>
                </div>
            `;
    });

    showResult("logsResult", logsHtml, "info");
  } catch (error) {
    showResult("logsResult", `❌ Error: ${error.message}`, "error");
  }
}

// Clear webhook logs
async function clearLogs() {
  const endpointId = document.getElementById("logsEndpointId").value.trim();

  if (!endpointId) {
    showResult("logsResult", "Please enter an endpoint ID", "error");
    return;
  }

  // Show confirmation dialog
  const confirmed = confirm(`Are you sure you want to clear ALL webhook logs for endpoint "${endpointId}"?\n\nThis action cannot be undone.`);

  if (!confirmed) {
    return;
  }

  showLoading("logsResult");

  try {
    const result = await makeRequest(`${API_BASE_URL}/webhooks/${endpointId}/logs`, {
      method: "DELETE",
    });

    const successMessage = `
      <h4>🗑️ Webhook logs cleared successfully!</h4>
      <div class="meta">
        <strong>Endpoint:</strong> ${result.endpointId}<br>
        <strong>Deleted Count:</strong> ${result.deletedCount} logs<br>
        <strong>Message:</strong> ${result.message}
      </div>
      <div class="data">
        <strong>Response:</strong><br>
        ${formatJSON(result)}
      </div>
    `;

    showResult("logsResult", successMessage, "success");
  } catch (error) {
    showResult("logsResult", `❌ Error: ${error.message}`, "error");
  }
}

// Generate endpoint ID
function generateEndpointId() {
  // Generate a random endpoint ID with format: prefix-timestamp-random
  const prefixes = ["webhook", "api", "hook", "endpoint", "test"];
  const prefix = prefixes[Math.floor(Math.random() * prefixes.length)];
  const timestamp = Date.now().toString(36);
  const random = Math.random().toString(36).substring(2, 8);

  const endpointId = `${prefix}-${timestamp}-${random}`;

  // Update both endpoint ID fields
  document.getElementById("endpointId").value = endpointId;
  document.getElementById("logsEndpointId").value = endpointId;

  // Show success message
  const successMessage = `
    <h4>🎯 Endpoint ID Generated!</h4>
    <div class="meta">
      <strong>Generated ID:</strong> ${endpointId}<br>
      <strong>Format:</strong> prefix-timestamp-random<br>
      <strong>Status:</strong> Ready to use
    </div>
    <div class="data">
      <strong>Details:</strong><br>
      • Prefix: ${prefix}<br>
      • Timestamp: ${new Date().toISOString()}<br>
      • Random: ${random}<br><br>
      <strong>Usage:</strong><br>
      This endpoint ID has been automatically filled in both the "Test Webhook" and "Webhook Logs" sections.
    </div>
  `;

  showResult("endpointIdResult", successMessage, "success");
}

// Toggle expandable sections
function toggleSection(sectionId) {
  const content = document.getElementById(sectionId);
  const header = content.previousElementSibling;
  const icon = header.querySelector(".toggle-icon");

  if (content.style.display === "none") {
    content.style.display = "block";
    icon.textContent = "▼";
    header.classList.remove("collapsed");
  } else {
    content.style.display = "none";
    icon.textContent = "▶";
    header.classList.add("collapsed");
  }
}

// Check health
async function checkHealth() {
  showLoading("healthResult");

  try {
    const result = await makeRequest(`${API_BASE_URL}/health`);

    const healthMessage = `
            <h4>✅ Service is healthy!</h4>
            <div class="meta">
                <strong>Status:</strong> ${result.status}<br>
                <strong>Service:</strong> ${result.service}<br>
                <strong>Message:</strong> ${result.message}
            </div>
            <div class="data">
                <strong>Response:</strong><br>
                ${formatJSON(result)}
            </div>
        `;

    showResult("healthResult", healthMessage, "success");
  } catch (error) {
    showResult("healthResult", `❌ Service is down: ${error.message}`, "error");
  }
}

// Show replay dialog
function showReplayDialog(webhookLogId, method) {
  const targetUrl = prompt(`Enter the target URL to replay this ${method} webhook:`, "https://your-endpoint.com/webhook");

  if (!targetUrl) {
    return; // User cancelled
  }

  const timeout = prompt("Enter timeout in seconds (default: 30):", "30");
  const timeoutValue = timeout ? parseInt(timeout) : 30;

  if (isNaN(timeoutValue) || timeoutValue <= 0) {
    alert("Invalid timeout value. Using default 30 seconds.");
    timeoutValue = 30;
  }

  replayWebhook(webhookLogId, targetUrl, timeoutValue);
}

// Replay webhook to external endpoint
async function replayWebhook(webhookLogId, targetUrl, timeout) {
  try {
    const result = await makeRequest(`${API_BASE_URL}/webhooks/replay/${webhookLogId}`, {
      method: "POST",
      body: JSON.stringify({
        target_url: targetUrl,
        timeout: timeout,
      }),
    });

    const replayResult = result.result;
    const statusIcon = replayResult.success ? "✅" : "❌";
    const statusText = replayResult.success ? "Success" : "Failed";

    const replayMessage = `
      <h4>${statusIcon} Webhook Replayed!</h4>
      <div class="meta">
        <strong>Target URL:</strong> ${targetUrl}<br>
        <strong>Status:</strong> ${statusText}<br>
        <strong>Response Code:</strong> ${replayResult.status_code}<br>
        <strong>Duration:</strong> ${replayResult.duration}<br>
        <strong>Webhook ID:</strong> ${webhookLogId}
      </div>
      <div class="data">
        <strong>Response Headers:</strong><br>
        ${formatJSON(replayResult.headers)}<br><br>
        <strong>Response Body:</strong><br>
        ${replayResult.response_body || "No response body"}
      </div>
    `;

    // Show result in a temporary alert or create a dedicated result area
    alert(
      `Webhook replay ${statusText.toLowerCase()}!\n\nStatus Code: ${replayResult.status_code}\nDuration: ${
        replayResult.duration
      }\n\nCheck the console for full details.`
    );
    console.log("Replay Result:", result);
  } catch (error) {
    alert(`❌ Replay failed: ${error.message}`);
    console.error("Replay Error:", error);
  }
}

// Auto-load logs when page loads
document.addEventListener("DOMContentLoaded", function () {
  // Check health on page load
  checkHealth();

  // Add keyboard shortcuts
  document.addEventListener("keydown", function (e) {
    // Ctrl+Enter to send webhook
    if (e.ctrlKey && e.key === "Enter") {
      if (document.activeElement.id === "webhookData") {
        sendWebhook();
      }
    }

    // Ctrl+L to load logs
    if (e.ctrlKey && e.key === "l") {
      e.preventDefault();
      loadLogs();
    }
  });

  // Add input validation
  document.getElementById("webhookData").addEventListener("input", function () {
    const data = this.value.trim();
    if (data) {
      try {
        JSON.parse(data);
        this.style.borderColor = "#28a745";
      } catch (error) {
        this.style.borderColor = "#dc3545";
      }
    } else {
      this.style.borderColor = "#ecf0f1";
    }
  });
});
