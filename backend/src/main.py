from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from routers import routers
from database.core import create_tables
from authx_extra.oauth2 import MiddlewareOauth2
import uvicorn
from config import settings

app = FastAPI(
    title="Sentimenta",
    description="Sentimenta is a mood tracker app",
    version="0.1.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
    )
app.add_middleware(
    MiddlewareOauth2,
    providers={
        'google': {
            'keys': 'https://www.googleapis.com/oauth2/v3/certs',
            'issuer': 'https://accounts.google.com',
            'audience': f'{settings.OAUTH2_GOOGLE_CLIENT_ID}.apps.googleusercontent.com',
        }
    },
    public_paths={
        "/",
        "/docs",
        "/favicon.ico",
        "/openapi.json",
        "/api/auth/register",
        "/api/auth/login",

        },  # Пути, не требующие аутентификации
)

[app.include_router(router) for router in routers]

if __name__ == "__main__":
    create_tables(True)
    uvicorn.run(app, host="0.0.0.0", port=8000)
