import os


POSTGRES_HOST = "bizkit-database"
POSTGRES_DB = os.environ.get("POSTGRES_DB")
POSTGRES_USER = os.environ.get("POSTGRES_USER")
POSTGRES_PASSWORD = os.environ.get("POSTGRES_PASSWORD")


event_types = [
    "Выставки",
    "Конференции",
    "Круглые столы",
    "Форумы",
    "Семинары",
    "Другое",
]

cache_timeout = 60 * 60
