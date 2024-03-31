The web app uses Single Store and Sharded MySQL for data storage.

1. Start SingleStore
    ```bash
    sudo docker pull ghcr.io/singlestore-labs/singlestoredb-dev:latest
    ```
    ```bash
    sudo docker run  -d --name singlestoredb-dev  -e ROOT_PASSWORD="root" \
        -e SINGLESTORE_VERSION="7.8" \ 
        --hostname singlestore-local \ 
        -e SINGLESTORE_LICENSE="<LICENSE>"  \ 
        -p 3306:3306 -p 8080:8080 -p 9000:9000 \ 
        -v my_cool_volume:/home/prakash/GolandProjects/singlestore-data  ghcr.io/singlestore-labs/singlestoredb-dev
    ```

    To connect to SingleStore
    ```bash
    sudo docker exec -it singlestoredb-dev singlestore -p
     ```
2. Start MySQL Shards

    Start with 2 shards and scale later.
   
    ```bash
    sudo docker pull mysql:8.0.32
    ```
    
    ```bash
    sudo docker run --name mysql1 -v ~/mysql-data:/home/prakash/GolandProjects/mysql-data-1  \
        -e MYSQL_ROOT_PASSWORD=root -p 3310:3306 -d mysql:8.0.32
    ```
    
    ```bash
    sudo docker run --name mysql2 -v ~/mysql-data:/home/prakash/GolandProjects/mysql-data-2 \
        -e MYSQL_ROOT_PASSWORD=root -p 3311:3306 -d mysql:8.0.32
    ```
    
    To connect to MySQL shards
    
    ```bash
    sudo mysql -h 127.0.0.1 -P 3310 -u root -p
    ```
    
    ```bash
    sudo mysql -h 127.0.0.1 -P 3311 -u root -p
    ```
    

3. Create tables   
    3.1 Create tables from idgenms in SingleStore   
    3.2 Create tables from database-clustermgt-ms in SingleStore   
    3.3 Create tables from main-url-shortener-ms in all MySQL Shards and domain_shortening_counts table in SingleStore   


4. Start idgenms, database-clustermgt-ms and main-url-shortener-ms in the same order as mentioned below   

    4.1 Start idgenms
    
    ```bash
    sudo docker pull prakashp92/idgenms:latest
    ```

    Execute idgenms

    ```bash
    sudo docker run -it -p 3001:3001 --name idgenms1 prakashp92/idgenms:latest
    ```
   
   
    
    4.2 Start database-clustermgt-ms
    
    ```bash
    sudo docker pull prakashp92/database-clustermgt-ms:latest
    ```

    Execute database-clustermgt-ms

    ```bash
    sudo docker run -it -p 3002:3002 --name clustermgtms1 prakashp92/database-clustermgt-ms:latest

    ```
   
   
    
    4.3 Start main-url-shortener-ms
    
    ```bash  
    sudo docker pull prakashp92/main-url-shortener-ms:latest
    ```    

    Execute main-url-shortener-ms
   
    ```bash  
    sudo docker run -it -p 3000:3000 --name main_url_shortener_ms1 prakashp92/main-url-shortener-ms:latest
    ````   
   
