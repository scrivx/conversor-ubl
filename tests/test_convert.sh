curl -X POST http://localhost:8080/convert `
  -H "Content-Type: application/json" `
  -d '{
    "document_type": "invoice",
    "data": {
      "id": "F001-123",
      "emisor_ruc": "20123456789",
      "emisor_nombre": "Mi Empresa SAC",
      "emisor_razon": "Mi Empresa S.A.C.",
      "cliente_ruc": "20654321098",
      "cliente_razon": "Cliente S.A.",
      "item_nombre": "Producto de prueba",
      "total": 150.75
    }
  }'
