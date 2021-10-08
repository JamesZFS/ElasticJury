# ElasticJury 0.1.1


![image-20200601010515162](./assets/screenshot1.png)


An archetypical search engine for judiciary cases.

By Fengshi Zheng and Chenggang Zhao

Please take a look at our [presentation slides](./docs/presentation.key) or [report](./docs/report.pdf) :)

## Build & Deploy

- To preprocess the dataset, build retrieve indexes and store to database, see `preprocessor/main.py`

- To start the search engine backend, run `PORT=8000; PASSWORD=${your db passwd} go run .` in the root directory.

- To start the frontend ui server in dev mode, run `npm run serve` in `./ui/`

- To build the frontend ui for production, run `npm run build` in `./ui`
    > Note: to deploy the static html files, you need a frontend server to serve `./ui/dist/`
    > 
    > For example, a sample nginx configuration file is located at `./ui/nginx.conf`, go check it out!

Have fun!

![image-20211008144203212](assets/image-20211008144203212.png)

![image-20211008144137959](assets/image-20211008144137959.png)
