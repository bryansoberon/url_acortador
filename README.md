# Smart URL Shortener

Este proyecto implementa un **acortador de URLs inteligente y resistente**, utilizando solo herramientas estándar de Go. Cumple con los requisitos de concurrencia, modularidad, validaciones y uso eficiente de memoria.

---

## 🚀 Endpoints

### 1. `POST /shorten`

Genera una URL corta para una URL larga.


ALUMNOS: 
BERNABE BRYAN SOBERON QUINTANA
FRANK GIANPIER ZEÑA VASQUEZ
DOCENTE:
ALEX JAVIER VILLEGAS LAINAS
TLP "A"

#### Solicitud

```json
{
  "long_url": "https://example.com/recurso/largo?query=param"
}
Respuesta
{
  "short_url": "http://localhost:8080/abc123"
}
```
