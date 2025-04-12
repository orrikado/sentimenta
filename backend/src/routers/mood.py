from fastapi import APIRouter, Depends, HTTPException
from database.schemas.mood import MoodAddSchema, MoodUpdateSchema
from database.queries.mood import add_mood, get_moods, update_mood
from jwt_token import security
from authx import TokenPayload

router = APIRouter(tags=["Mood"])


@router.post("/api/moods/add", dependencies=[Depends(security.access_token_required)])
async def add_mood_route(
    mood: MoodAddSchema, user: TokenPayload = Depends(security.access_token_required)
):
    try:
        user_id = int(user.sub)
        await add_mood(user_id, mood)
    except Exception as e:
        raise HTTPException(status_code=404, detail=f"{e}")


@router.get("/api/moods/get", dependencies=[Depends(security.access_token_required)])
async def get_moods_route(user: TokenPayload = Depends(security.access_token_required)):
    user_id = int(user.sub)
    moods = await get_moods(user_id)
    if moods:
        return moods
    else:
        raise HTTPException(status_code=404, detail="Mood not found")


@router.put("/api/moods/update", dependencies=[Depends(security.access_token_required)])
async def update_mood_route(
    mood: MoodUpdateSchema, user: TokenPayload = Depends(security.access_token_required)
):
    user_id = int(user.sub)
    await update_mood(user_id, mood)
