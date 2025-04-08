from src.database.db_setup import sync_engine
from src.database.db_setup import Base

async def create_tables(drop_tables: bool = False):
    if drop_tables:
        Base.metadata.drop_all(sync_engine)
    Base.metadata.create_all(sync_engine)
