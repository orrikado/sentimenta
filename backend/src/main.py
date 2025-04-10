from fastapi import FastAPI
from routers import routers
from database.core import create_tables
import uvicorn

app = FastAPI()

[app.include_router(router) for router in routers]

if __name__ == "__main__":
    create_tables(True)
    uvicorn.run(app, host="0.0.0.0", port=8000)
