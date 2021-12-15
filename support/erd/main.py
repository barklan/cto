from typing import Optional


from fastapi import FastAPI, File, UploadFile
from eralchemy import render_er
from fastapi.staticfiles import StaticFiles
from fastapi.responses import HTMLResponse

import random, string

def randomword(length):
   letters = string.ascii_lowercase
   return ''.join(random.choice(letters) for i in range(length))

app = FastAPI(debug=True)


@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.post("/uploadfile/")
async def create_upload_file(file: UploadFile = File(...)):
    contents = await file.read()

    randname = randomword(12)
    with open(f'/app/media/{randname}.db', 'wb') as outfile:
        outfile.write(contents)

    render_er(f"sqlite:////app/media/{randname}.db", f'{randname}.png')
    return {"filename": file.filename}


app.mount("/static", StaticFiles(directory="./"), name="static")
app.mount("/front", StaticFiles(directory="/app/front/"), name="front")
