import time
import logging
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse, Response
from prometheus_client import Counter, Histogram, Gauge, generate_latest, CONTENT_TYPE_LATEST
import uvicorn

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger(__name__)

app = FastAPI()

# Prometheus metrics setup
requests_total = Counter(
    "api_requests_total",
    "Total de requests recebidos",
    ["method", "route", "status"]
)

request_duration = Histogram(
    "api_request_duration_seconds",
    "Duração das requests em segundos",
    ["method", "route"]
)

requests_in_flight = Gauge(
    "api_requests_in_flight",
    "Requests sendo processadas no momento"
)

@app.middleware("http")
async def metrics_middleware(request: Request, call_next):
    route = request.url.path
    if not route:
        route = "unknown"
        
    requests_in_flight.inc()
    start_time = time.time()
    
    try:
        response = await call_next(request)
        status_code = str(response.status_code)
    except Exception as e:
        status_code = "500"
        requests_in_flight.dec()
        raise e

    duration = time.time() - start_time
    requests_in_flight.dec()
    
    # Record metrics matches exactly the Go middleware
    requests_total.labels(method=request.method, route=route, status=status_code).inc()
    request_duration.labels(method=request.method, route=route).observe(duration)
    
    return response

@app.on_event("startup")
async def startup_event():
    logger.info("Iniciando blunteds-devops-api...")
    logger.info("Server listening on :8080")

@app.get("/health")
async def health(request: Request):
    client_ip = request.client.host if request.client else "unknown"
    logger.info(f"health check requested from {client_ip}")
    return {"status": "ok"}

@app.get("/")
async def root(request: Request):
    client_ip = request.client.host if request.client else "unknown"
    logger.info(f"root requested from {client_ip}")
    return {"message": "Bem-vindo à DevOps API! v2"}

@app.get("/metrics")
async def metrics():
    return Response(generate_latest(), media_type=CONTENT_TYPE_LATEST)

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8080)
