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

# PARA VENTA GRAVADA Y EXONERADA
curl -X POST http://localhost:8080/convert `
  -H "Content-Type: application/json" `
  -d '{
  "document_type": "invoice",
  "data": {
    "id": "F001-200",
    "emisor_ruc": "20123456789",
    "emisor_nombre": "Empresa SAC",
    "emisor_razon": "Empresa S.A.C.",
    "cliente_ruc": "20654321098",
    "cliente_razon": "Cliente SAC",
    "item_nombre": "Servicio A",
    "venta_gravada": 150.75,
    "venta_exonerada": 50.00,
    "total": 200.75
  }
}'

# TRANSFERENCIA GRATUITA
curl -X POST http://localhost:8080/convert `
  -H "Content-Type: application/json" `
  -d '{
  "document_type": "invoice",
  "data": {
    "id": "F001-201",
    "emisor_ruc": "20123456789",
    "emisor_nombre": "Empresa SAC",
    "emisor_razon": "Empresa S.A.C.",
    "cliente_ruc": "20654321098",
    "cliente_razon": "Cliente SAC",
    "item_nombre": "Muestra Gratis",
    "gratuita": 80.00
  }
}'

# PARA VENTA CON BONIFICACION
curl -X POST http://localhost:8080/convert `
  -H "Content-Type: application/json" `
  -d '{
  "document_type": "invoice",
  "data": {
    "id": "F001-202",
    "emisor_ruc": "20123456789",
    "emisor_nombre": "Empresa SAC",
    "emisor_razon": "Empresa S.A.C.",
    "cliente_ruc": "20654321098",
    "cliente_razon": "Cliente SAC",
    "item_nombre": "Producto Bonificado",
    "venta_gravada": 200.00,
    "bonificacion": 20.00
  }
}'

# PARA VENTA CON PERCEPCION
curl -X POST http://localhost:8080/convert `
  -H "Content-Type: application/json" `
  -d '{
  "document_type": "invoice",
  "data": {
    "id": "F001-203",
    "emisor_ruc": "20123456789",
    "emisor_nombre": "Empresa SAC",
    "emisor_razon": "Empresa S.A.C.",
    "cliente_ruc": "20654321098",
    "cliente_razon": "Cliente SAC",
    "item_nombre": "Bienes con Percepción",
    "venta_gravada": 300.00,
    "percepcion": 15.00
  }
}'

# PARA VENTA AL CREDITO
curl -X POST http://localhost:8080/convert `
  -H "Content-Type: application/json" `
  -d '{
  "document_type": "invoice",
  "data": {
    "id": "F001-204",
    "emisor_ruc": "20123456789",
    "emisor_nombre": "Empresa SAC",
    "emisor_razon": "Empresa S.A.C.",
    "cliente_ruc": "20654321098",
    "cliente_razon": "Cliente SAC",
    "item_nombre": "Servicio Financiado",
    "venta_gravada": 500.00,
    "credito": true,
    "plazo": "30 días"
  }
}'