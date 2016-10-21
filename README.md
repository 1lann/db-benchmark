# DB Benchmark
DB Benchmark is a framework to benchmark database performance with Go.

## Notices
I've made my best effort to write queries as idiomatically and realistically as possible. If you feel like that are issues with the queries used please don't hesitate to post a issue or pull request.

Here are some things to take note of with these database configurations:
- MySQL is running with default settings out of the box.
- RethinkDB is running with default settings out of the box.
- RethinkDB writes queries are performed with soft durability.
- 80 connections are used in the connection pools for both databases.

## Disclaimer
Note that purely the drivers of each database would have a significant effect on the results, as the difference between MySQL and RethinkDB are fractions of milliseconds per query. This benchmark is could be more reflective of Go driver performance rather than actual database performance.

I am in no way qualified to benchmark databases and draw conclusions from these results. This was really just an experiment and for fun. Please make your own judgements before deciding on the database appropriate for you, performance isn't everything. Remember that this benchmark does not take into account of different features, consistency, scalability and durability.

## Results
[results.txt](/results.txt)

## License
[MIT](/LICENSE)
