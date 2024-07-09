window.onload = function() {
    // Build a system
    const ui = SwaggerUIBundle({
      url: "/v1/swagger-v1.yaml",
      dom_id: '#swagger-ui',
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIStandalonePreset
      ],
      layout: "StandaloneLayout"
    })
  
    window.ui = ui
  }
  