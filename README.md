# site-analyzer

To run via Docker

    docker build -t site-analyzer .
    docker run -v /data:/data -t site-analyzer -url https://nytimes.com
