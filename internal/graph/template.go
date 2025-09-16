package graph

// Build with AI ;)
// Because i don't like dev frontend
func (g *Generator) getHTMLTemplate() string {
	return `<!DOCTYPE html>
       <html>
       <head>
          <meta charset="utf-8">
          <title>PHP SPX Graph Visualizer</title>
          <style>
             * {
                margin: 0;
                padding: 0;
                box-sizing: border-box;
             }
             
             body {
                font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
                background: #f8f9fa;
                color: #333;
             }
             
             .header {
                background: #ffffff;
                border-bottom: 1px solid #e1e5e9;
                padding: 16px 24px;
             }
             
             .header h1 {
                font-size: 24px;
                font-weight: 600;
                color: #1a1a1a;
             }
             
             .header p {
                font-size: 14px;
                color: #6c757d;
                margin-top: 4px;
             }
             
             .controls {
                background: #ffffff;
                border-bottom: 1px solid #e1e5e9;
                padding: 12px 24px;
                display: flex;
                align-items: center;
                gap: 16px;
             }
             
             .stats {
                display: flex;
                gap: 16px;
             }
             
             .stat-item {
                background: #f8f9fa;
                padding: 8px 12px;
                border: 1px solid #e1e5e9;
                font-size: 14px;
             }
             
             .stat-value {
                font-weight: 600;
                color: #495057;
             }
             
             .zoom-controls {
                display: flex;
                gap: 8px;
                margin-left: auto;
             }
             
             .btn {
                background: #ffffff;
                border: 1px solid #d1d5db;
                padding: 6px 12px;
                cursor: pointer;
                font-size: 14px;
                transition: background-color 0.2s;
             }
             
             .btn:hover {
                background: #f3f4f6;
             }
             
             .btn:active {
                background: #e5e7eb;
             }
             
             .graph-container {
                position: relative;
                background: #ffffff;
                height: calc(100vh - 140px);
                overflow: hidden;
                border-top: 1px solid #e1e5e9;
             }
             
             .graph-viewport {
                width: 100%;
                height: 100%;
                overflow: auto;
                cursor: grab;
             }
             
             .graph-viewport:active {
                cursor: grabbing;
             }
             
             .graph-content {
                transform-origin: 0 0;
                transition: transform 0.1s ease-out;
                min-width: 100%;
                min-height: 100%;
             }
             
             svg {
                display: block;
                max-width: none;
                height: auto;
             }
          </style>
       </head>
       <body>
          <div class="header">
             <h1>SPX Profile Graph</h1>
             <p>Call graph with profiling data</p>
          </div>
          
          <div class="controls">
             <div class="stats">
                <div class="stat-item">
                   <span class="stat-value">{{.NodeCount}}</span> Functions
                </div>
                <div class="stat-item">
                   <span class="stat-value">{{.EdgeCount}}</span> Edges
                </div>
                <div class="stat-item">
                   <span class="stat-value">{{.TotalCalls}}</span> Total Calls
                </div>
             </div>
             
             <div class="zoom-controls">
                <button class="btn" onclick="zoomIn()">Zoom In</button>
                <button class="btn" onclick="zoomOut()">Zoom Out</button>
                <button class="btn" onclick="resetZoom()">Reset</button>
                <button class="btn" onclick="fitToScreen()">Fit</button>
             </div>
          </div>
          
          <div class="graph-container">
             <div class="graph-viewport" id="viewport">
                <div class="graph-content" id="content">
                   {{.SVG}}
                </div>
             </div>
          </div>
          
          <script>
             let scale = 1;
             let panX = 0;
             let panY = 0;
             let isPanning = false;
             let lastX = 0;
             let lastY = 0;
             
             const viewport = document.getElementById('viewport');
             const content = document.getElementById('content');
             
             function zoomIn() {
                scale *= 1.2;
                updateTransform();
             }
             
             function zoomOut() {
                scale /= 1.2;
                updateTransform();
             }
             
             function resetZoom() {
                scale = 1;
                panX = 0;
                panY = 0;
                updateTransform();
             }
             
             function fitToScreen() {
                const svg = content.querySelector('svg');
                if (!svg) return;
                
                const svgRect = svg.getBoundingClientRect();
                const containerRect = viewport.getBoundingClientRect();
                
                const scaleX = containerRect.width / svgRect.width;
                const scaleY = containerRect.height / svgRect.height;
                scale = Math.min(scaleX, scaleY, 1) * 0.9;
                
                panX = 0;
                panY = 0;
                updateTransform();
             }
             
             function updateTransform() {
                content.style.transform = 'translate(' + panX + 'px, ' + panY + 'px) scale(' + scale + ')';
             }
             
             viewport.addEventListener('wheel', function(e) {
                e.preventDefault();
                const delta = e.deltaY;
                const zoomFactor = delta > 0 ? 0.9 : 1.1;
                
                const rect = viewport.getBoundingClientRect();
                const mouseX = e.clientX - rect.left;
                const mouseY = e.clientY - rect.top;
                
                const oldScale = scale;
                scale *= zoomFactor;
                
                panX = mouseX - (mouseX - panX) * (scale / oldScale);
                panY = mouseX - (mouseY - panY) * (scale / oldScale);
                
                updateTransform();
             });
             
             viewport.addEventListener('mousedown', function(e) {
                isPanning = true;
                lastX = e.clientX;
                lastY = e.clientY;
                viewport.style.cursor = 'grabbing';
             });
             
             document.addEventListener('mousemove', function(e) {
                if (!isPanning) return;
                
                const deltaX = e.clientX - lastX;
                const deltaY = e.clientY - lastY;
                
                panX += deltaX;
                panY += deltaY;
                
                lastX = e.clientX;
                lastY = e.clientY;
                
                updateTransform();
             });
             
             document.addEventListener('mouseup', function() {
                isPanning = false;
                viewport.style.cursor = 'grab';
             });
             
             viewport.addEventListener('touchstart', function(e) {
                if (e.touches.length === 1) {
                   isPanning = true;
                   lastX = e.touches[0].clientX;
                   lastY = e.touches[0].clientY;
                }
             });
             
             viewport.addEventListener('touchmove', function(e) {
                e.preventDefault();
                
                if (e.touches.length === 1 && isPanning) {
                   const deltaX = e.touches[0].clientX - lastX;
                   const deltaY = e.touches[0].clientY - lastY;
                   
                   panX += deltaX;
                   panY += deltaY;
                   
                   lastX = e.touches[0].clientX;
                   lastY = e.touches[0].clientY;
                   
                   updateTransform();
                }
             });
             
             viewport.addEventListener('touchend', function() {
                isPanning = false;
             });
             
             window.addEventListener('load', function() {
                resetZoom();
             });
          </script>
       </body>
       </html>`
}
