# golang-api

## Simple API with Go, RabbitMQ, MySQL and ElasticSearch

How to run
- clone this repository
- go to cloned dir `cd golang-api`
- run docker `docker-compose up`
- wait until all service run
- run your POST and GET request

Progress List:
- [x] Queue Message when post to RabbitMQ
- [x] Worker to save queued data to MySQL and ElasticSearch
- [x] Implement goroutine when getting data and running worker
- [x] Cache result data
- [ ] Testing
- [ ] Make application wait for all service in docker, right now. When all services not run and we run the API request the application will exit. So we need to re run the application again with `docker-compose up app`

## Request & Response Examples

### API Resources

  - [GET /news/?page=1](#get-news)
  - [POST /news](#post-news)

### GET /magazines

Example: http://localhost:3030/news

Request Param (Query String):

    [
        "page": 1
    ]

Response body:

    {
        "data":[
            {
                "id": 2,
                "author": "John Doe",
                "body": "New Lorem Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam malesuada accumsan nisl nec bibendum. Nam nec nibh commodo, porta quam non, aliquam dolor. Vivamus tincidunt scelerisque purus vel laoreet. Quisque eleifend dolor ut ligula vulputate, quis convallis nulla ullamcorper. Nulla ultricies nisl nec mi semper, non porta elit facilisis. Aenean et laoreet diam, at rutrum velit. Praesent tempus purus vitae felis fringilla accumsan. Maecenas vulputate sed metus auctor lobortis. Sed venenatis neque blandit molestie pretium. Maecenas at justo sed risus sodales ultrices. Aliquam erat volutpat",
                "created": "2019-03-25T07:04:16Z"
            },
            {
                "id": 1,
                "author": "Puas Citro",
                "body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam malesuada accumsan nisl nec bibendum. Nam nec nibh commodo, porta quam non, aliquam dolor. Vivamus tincidunt scelerisque purus vel laoreet. Quisque eleifend dolor ut ligula vulputate, quis convallis nulla ullamcorper. Nulla ultricies nisl nec mi semper, non porta elit facilisis. Aenean et laoreet diam, at rutrum velit. Praesent tempus purus vitae felis fringilla accumsan. Maecenas vulputate sed metus auctor lobortis. Sed venenatis neque blandit molestie pretium. Maecenas at justo sed risus sodales ultrices. Aliquam erat volutpat",
                "created": "2019-03-25T07:04:16Z"
            }
        ],
        "limit": 10,
        "page": 1
    }


### GET /magazines/[id]

Example: http://localhost:3030/news

Request body:

    [
        {
            "author": "John Doe",
            "body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam eget ante ut augue scelerisque ornare. Aliquam tempus rhoncus quam vel luctus. Sed scelerisque fermentum fringilla. Suspendisse tincidunt nisl a metus feugiat vitae vestibulum enim vulputate. Quisque vehicula dictum elit, vitae cursus libero auctor sed. Vestibulum fermentum elementum nunc. Proin aliquam erat in turpis vehicula sit amet tristique lorem blandit. Nam augue est, bibendum et ultrices non, interdum in est. Quisque gravida orci lobortis... "
        }
    ]

Response body:

    {
        "message": "News queued for saving"
    }
