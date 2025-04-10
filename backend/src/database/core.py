from database.db_setup import sync_engine
from database.db_setup import Base

def create_tables(drop_tables: bool = False):
    if drop_tables:
        Base.metadata.drop_all(sync_engine)
    Base.metadata.create_all(sync_engine)
    print("Database tables created")
