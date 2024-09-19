import random
from uuid import UUID

import asyncpg
from aiocache import cached, Cache
from aiocache.serializers import PickleSerializer

from misc.config import (
    POSTGRES_DB,
    POSTGRES_USER,
    POSTGRES_PASSWORD,
    cache_timeout,
    event_types,
    POSTGRES_HOST,
)


custom_cached = cached(
    cache=Cache.MEMORY,
    serializer=PickleSerializer(),
    ttl=cache_timeout,
)


class Database:
    _connection: asyncpg.Connection | None = None

    @classmethod
    async def create(cls) -> None:
        cls._connection = await asyncpg.connect(
            host=POSTGRES_HOST,
            database=POSTGRES_DB,
            user=POSTGRES_USER,
            password=POSTGRES_PASSWORD,
        )

    @classmethod
    async def close(cls) -> None:
        await cls._connection.close()

    @classmethod
    @custom_cached
    async def get_unvisited_events(cls, user_uuid: UUID) -> dict[str, list[UUID]]:
        query = """
            SELECT e.type, ARRAY_AGG(e.uuid ORDER BY e.updated_at DESC) as events
            FROM events e
            LEFT JOIN events_users eu
            ON e.uuid = eu.event_uuid AND eu.user_uuid = $1
            WHERE eu.user_uuid IS NULL
            GROUP BY e.type
        """

        grouped_events = await cls._connection.fetch(query, str(user_uuid))

        return {event_type: [] for event_type in event_types} | {
            event_type_dict["type"]: event_type_dict["events"]
            for event_type_dict in grouped_events
        }

    @classmethod
    @custom_cached
    async def get_visit_statistics(cls, user_uuid: UUID) -> dict[str, float]:
        total_visited_count: int = await cls._connection.fetchval(
            "SELECT COUNT(*) FROM events_users WHERE user_uuid = $1", str(user_uuid)
        )

        query = """
            SELECT e.type, COUNT(*) as count
            FROM events e
            LEFT JOIN events_users eu
            ON e.uuid = eu.event_uuid
            WHERE eu.user_uuid = $1
            GROUP BY e.type
        """

        statistics = await cls._connection.fetch(query, str(user_uuid))

        return {event_type: 0 for event_type in event_types} | {
            event_type_dict["type"]: event_type_dict["count"] / total_visited_count
            for event_type_dict in statistics
        }

    @classmethod
    @custom_cached
    async def get_recommendations(cls, user_uuid: UUID) -> list[UUID]:
        unvisited_events = await cls.get_unvisited_events(user_uuid)
        named_weights = await cls.get_visit_statistics(user_uuid)

        if sum(named_weights.values()) == 0:
            result = [
                event
                for category in event_types
                for event in unvisited_events[category]
            ]

            random.shuffle(result)

            return result

        result = []

        for _ in range(sum(map(len, unvisited_events.values()))):
            weights = list(named_weights.values())

            category = random.choices(list(unvisited_events.keys()), weights)[0]

            result.append(unvisited_events[category].pop(0))

            if not unvisited_events[category]:
                named_weights[category] = 0

        return result
