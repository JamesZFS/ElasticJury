# ElasticJury

An archetypical search engine for judiciary cases.

By Lyric Zhao and James Zheng

## Build & Deploy

- To preprocess the dataset, build retrieve indexes and store to database, see `preprocessor/main.py`

- To start the search engine backend, run `PORT=8000; PASSWORD=${your db passwd} go run .` in the root directory.

- To start the frontend ui server in dev mode, run `npm run serve` in `./ui/`

- To build the frontend ui for production, run `npm run build` in `./ui`
    > Note: to deploy the static html files, you need a frontend server to serve `./ui/dist/`
    > 
    > For example, a sample nginx configuration file is located at `./ui/nginx.conf`, go check it out!

Have fun!
