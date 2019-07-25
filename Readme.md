##**Starting the server**  
In the current directory  

    ./start.sh  
 The server will listen on port 8001
    
##**APIS**

###***Get Products***

    GET - /api/v1/products

###***Index Product***

    POST - /api/v1/products
   ***POST BODY***   
   

     {"title":"shoe","brand":"roadster","price":1000,"stock":10}

###***Search***

    GET - /api/v1/products?q=\<search_terms\>

###***Filter***

    GET - /api/v1/products?filter=\<filter_field\>:\<filter_value\>

###***Sorting***

    GET - /api/v1/products?sort=\<sort_field\>

###***Pagination***

    GET - /api/v1/products?page=\<int\>&size=\<int\>

##***View Logs***

    docker-compose logs -f goGFG
