from contextlib import asynccontextmanager, AbstractAsyncContextManager
from uuid import UUID

from fastapi import FastAPI, Body
from fastapi.middleware.cors import CORSMiddleware

from logic.database import Database


# noinspection PyUnusedLocal
@asynccontextmanager
async def lifespan(app: FastAPI) -> AbstractAsyncContextManager[None]:
    await Database.create()

    yield

    await Database.close()


app = FastAPI(lifespan=lifespan)

# noinspection PyTypeChecker
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "localhost",
    ],
    allow_credentials=True,
    allow_methods=["GET", "POST"],
    allow_headers=["*"],
)


@app.post("/recommendations/recommend")
async def recommend(
    user_uuid: UUID = Body(embed=True),
    limit: int = Body(embed=True, default=10),
    offset: int = Body(embed=True, default=0),
) -> list[UUID]:
    all_recommendations = await Database.get_recommendations(user_uuid)

    return all_recommendations[offset : offset + limit]
