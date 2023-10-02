# Scraper
`scraper` is a tool for scraping websites.

## Design
`scraper` scrapes website by recursively traversing the links found on each page. Each recursive call is handled by an asynchronous Go func, this allows `scraper` to be very fast while the implementation is kept fairly simple.

I designed `scraper` with unit testing in mind, I was important to me that certain operations could be mocked in order to test the actual logic. this is why the `Scraper` type have references to `http.Fetcher`, `links.Parser`, `filesystem.FileWriter`. By relying on these interfaces in the  `Scraper.Scrape`-function, it becomes easy to write unit tests covering the entire function, including all possible error scenarios. 

## How to build
`scraper` can be built using the following command
```
make build
```
This will test and build the code. The output is `scaper.exe`.

## How to run
Running `scraper.exe -h` reveals the following help message
```
> scraper.exe -h
Usage of scraper.exe:
  -o string
        Output directory
  -s string
        The site to scrape
```

To scrape http://books.toscrape.com, run
```
scraper.exe -s http://books.toscrape.com -o /path/to/output/folder
```

## How to run the unit tests
To run the unit tests, use  
```
make test
```

## How it works
The main logic of `scraper` lives in the `scraping`-package, the orchestrator of the scraping operations is the `Scraper` type.

### `Scraper`
The `Scraper`-type holds the following attributes
* `fetcher` - For fetching websites
* `fileWriter` - For writing files to the file system
* `linksParser` - For parsing links on a website
* `baseURL` - The root website URL to be scraped
* `outputDirectory` - The location on disk where the files should be downloaded
* `claimer` - `claimer` keeps track of all links that have been claimed in order to avoid scraping a link more than once
* `counter` - `counter` keeps track of the amount of downloaded assets

The `Scraper`-type contains a single exported function `Scraper.Scrape` that starts the scraping by calling the private recursive function `Scraper.scrape` with an empty string as the argument. `Scraper.scrape` does the following things:
1. Downloads the url is was passed (a relative url)
2. Writes the downloaded asset to disk
3. Increases the `counter`
4. If the url did not end in `.html` or `.htm` it returns, otherwise it continues
5. Parses the links on the downloaded page
6. Iterates over the parsed links
      1. Tries to claim the link, if it is already claimed, the iteration continues
      2. Spins up a Go func where `Scraper.scrape` is called with the new link as the argument
7. Wait for all Go funcs to finish or an error to occur and returns to the caller

## Limitations
`scraper` currently only handles relative URLs, adding support for absolute URLs would be trivial

## Potential caveats
`scraper` works perfectly fine when scraping http://books.toscrape.com but I suspect it might struggle with larger websites due to potentially hitting the max recursion depth or running out of Go funcs if there's a large amount of links. The recursion depth problem could be solved by implementing a depth system where the recursive function halts after hitting a set depth and then start back up with the links that were found on the last call. Running out of Go functions could be solved by implementing a worker pool.



