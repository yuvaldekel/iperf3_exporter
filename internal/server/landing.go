// Copyright 2026 Yuval Dekel
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

const LandingPageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>iPerf3 Exporter</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
            padding: 0;
            margin: 0;
            line-height: 1.5;
        }
        header {
            background-color: #1a73e8;
            color: white;
            padding: 1rem;
            margin-bottom: 1rem;
        }
        header h1 {
            margin: 0;
            font-size: 1.5rem;
        }
        .container {
            padding: 0 1rem;
            max-width: 1200px;
            margin: 0 auto;
        }
        .links {
            margin-bottom: 2rem;
        }
        .links a {
            color: #1a73e8;
            text-decoration: none;
            font-weight: bold;
        }
        .links a:hover {
            text-decoration: underline;
        }

        table {
            border-collapse: collapse;
            width: 100%%;
        }
        th, td {
            text-align: left;
            padding: 8px;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #f2f2f2;
        }
        pre {
            background-color: #f5f5f5;
            padding: 10px;
            border-radius: 5px;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <header>
        <div class="container">
            <h1>iPerf3 Exporter</h1>
        </div>
    </header>

    <div class="container">
        <p>The iPerf3 exporter allows iPerf3 probing of endpoints for Prometheus monitoring.</p>
        
        <div class="links">
            <a href="%s">Metrics</a>
        </div>

        <p>Version: %s</p>

        <h2>Quick Start</h2>
        <p>To probe a target:</p>
        <pre><a href="%s?target=example.com">%s?target=example.com</a></pre>

        <h2>Probe Parameters</h2>
        <table>
            <tr>
                <th>Parameter</th>
                <th>Description</th>
                <th>Default</th>
            </tr>
            <tr>
                <td>target</td>
                <td>Target host to probe (required)</td>
                <td>-</td>
            </tr>
            <tr>
                <td>port</td>
                <td>Port that the target iperf3 server is listening on</td>
                <td>5201</td>
            </tr>
            <tr>
                <td>reverse_mode</td>
                <td>Run iperf3 in reverse mode (server sends, client receives)</td>
                <td>false</td>
            </tr>
            <tr>
                <td>protocol</td>
                <td>Run iperf3 in UDP or TCP protocol</td>
                <td>tcp</td>
            </tr>
            <tr>
                <td>bitrate</td>
                <td>Target bitrate in bits/sec (format: #[KMG][/#]). For UDP mode, iperf3 defaults to 1 Mbit/sec if not specified.</td>
                <td>-</td>
            </tr>
            <tr>
                <td>period</td>
                <td>Duration of the iperf3 test</td>
                <td>5s</td>
            </tr>
        </table>

        <h2>Prometheus Configuration Example</h2>
        <pre>
scrape_configs:
  - job_name: 'iperf3'
    metrics_path: %s
    static_configs:
      - targets:
        - foo.server
        - bar.server
    params:
      port: ['5201']
      # Optional: enable reverse mode
      # reverse_mode: ['true']
      # Optional: protocol to use tcp or udp default to tcp
      # protocol: ['tcp']
      # Optional: set bitrate limit
      # bitrate: ['100M']
      # Optional: set test period
      # period: ['10s']
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9579  # The iPerf3 exporter's real hostname:port.
        </pre>
        </div>
</body>
</html>
`