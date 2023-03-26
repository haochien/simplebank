# external lib required by this project:
1. operate db migration (for windows system):
    ```bash
    $ scoop install migrate
    ```

2. lib/pq: to let go sql to talk to specific db engine
    ```bash
    go get github.com/lib/pq
    ```

3. testify: testing framework
    ```bash
    go get github.com/stretchr/testify
    ```

4. Golang Gin (Web framework)
    ```bash
    go get -u github.com/gin-gonic/gin
    ```

5. Golang Viper (env config )
    ```bash
    go get github.com/spf13/viper
    ```

5. Golang GoMock (mock DB)
    ```bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

6. uuid for Golang:
    ```bash
    go get github.com/google/uuid
    ```


7. Golang JWT:
    ```bash
    go get github.com/dgrijalva/jwt-go
    ```

8. Golang Paseto:
    ```bash
    go get github.com/o1egl/paseto
    ```

# Set up and connect postgres from docker:
0. pull docker image:
    ```docker
    docker pull postgres:15-alpine
    ```


1. connect to postgres:
    ```docker
    docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine
    ```
2. create/drop db
    ```docker
    docker exec -it postgres15 createdb --username=root --owner=root simple_bank

    docker exec -it postgres15 dropdb simple_bank
    ```

3. others:
    ```docker
    # stop db container
    docker stop postgres15

    # start db container
    docker start postgres15

    # remove current db container
    docker rm postgres15

    # log db 
    docker log postgres15

    # execute the postgres image as user root
    docker exec -it postgres15 psql -U root
    select now()
    ```


# Prepare migration:
1. init migration folder:
    ```bash
    migrate create -ext sql -dir db/migration -seq init_schema
    ```

2. run migration script:
    ```bash
    migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

    migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
    ```

3. we can continue on create other migration scripts:
    ```bash
    migrate create -ext sql -dir db/migration -seq second_script
    ```


# How to set up go module and dependency:
1. init module:
    ```bash
    go mod init github.com/haochien/simplebank
    ```
2. create/update dependency:
    ```bash
    go mod tidy
    ```


# How to use sqlc in go project:
0. download sqlc:
    ```docker
    docker pull kjconroy/sqlc
    ```
1. init sqlc: create sqlc.yaml. Can be done via:
    ```docker
    docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc init
    ```
2. configure content in sqlc.yaml
3. prepare sql files under query folder
4. generate go files sqlc functions via:
    ```docker
    docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc generate
    ```

5. create transaction maintain file: store.go 


# How to set up go Mock DB environment:
0. download gomock:
    ```bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```
 

1. all function under store struct should be map to store interface.
   we can manually add all those function to new store interface created in the next step.
   But easier way would be add following parameter in sqlc.yaml:
   ```yaml
   emit_interface: true
   ```
   and regenerate sqlc
    ```bash
    docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc generate

    # querier.go will then be created 
    ```

    

2. create a store interface to differentiate the real store struct
    ```go
    // in store.go file: 

    // Store defines all functions to execute db queries and transactions
    type Store interface {
        Querier
        TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
    }

    // SQLStore provides all functions to execute SQL queries and transactions
    type SQLStore struct {
        db *sql.DB
        *Queries
    }

    ```

3. create a mock folder under db folder, and run following command:
    ```bash
    mockgen -package mockdb -destination db/mock/store.go github.com/haochien/simplebank/db/sqlc Store

    # store.go will then created under mock folder
    ```


# How to build docker file:
1. create dockerfile

2. build the image
    ```docker
    docker build -t simplebank:latest .

    # in case rerun and need to delete old image (c8f6cb16c708 is image id):
    docker rmi c8f6cb16c708 
    ```

3. run created image
    ```docker
    # --name: specify the name of container we are going to create
    docker run --name simplebank -p 8080:8080 simplebank:latest

    # to run not in debug mode:
    docker run --name simplebank -p 8080:8080 -e GIN_MODE=release simplebank:latest

    # to deal with issue that default ip address of postgres and our simplebank is different:
    docker run --name simplebank -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@172.17.0.2:5432/simple_bank?sslmode=disable" simplebank:latest
    # note: you can use following command to check the ip address
    docker container inspect simplebank
    docker container inspect postgres15  

    # better solution to deal with different ip in two container: create your own network, instead of default bridge network
    docker network create bank-network
    docker network connect bank-network postgres15
    docker network inspect bank-network
    docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres15:5432/simple_bank?sslmode=disable" simplebank:latest


    # in case error occurs and need to rebuild:
    docker rm simplebank
    docker rmi c8f6cb16c708 
    docker build -t simplebank:latest .
    ```

4. create start.sh to able to run migration up command:

    remember to update migrate commands in dockerfile to enable migrate binary download and usage.

    also need update dockerfile to copy start.sh into image
    
    add ENTRYPOINT [ "/app/start.sh" ] in dockerfile to make start.sh as start point of the command

5. create docker-compose.yaml file to auto run step 3 commands:
    ```docker
    # to run docker compose file:
    # after run this, it will start from build image, and then connect to db container and access to api container in parallel
    docker compose up 

    # if need to revert, then run:
    docker compose down
    docker rmi simplebank-api
    ```

    Note that error might occur since db container and api container are build in parallel.

    The api should wait until db is set up (because api relies on db connection)

    Thus, depends_on should be setup in docker-compose.yaml.
    
    But depends_on does not wait for db and redis to be “ready” before starting web, healthcheck should be added to provide wait feature




# How to deploy to AWS:
1. auto build and push image to AWS ECR with github action:
    
    a. create 1 Repositories in Amazon ECR

    b. create 1 user in IAM (in this project example, remember to create Access key after user is created)

    c. create yaml for deploy action (in config, deploy.yml only runs when push to master; test.yml runs when raise PR or merge to master). check iin github market place for instruction: https://github.com/marketplace/actions/amazon-ecr-login-action-for-github-actions

    d. in github -> setting (under security -> secrets and variables -> Actions), configure deploy.yml  secret variables in Repository secrets


2. setup aws rds:

    a. create database in aws rds (get host(Endpoint), user, password, database, port)

    b. try to connect via TablePlus

    c. update migrate command in makefile. Migrate all db change to the database in aws


3. setup AWS Secrets Manager to mange secret keys and env variables:

    a. generate stronger secret keys (for example, to get 32-character TOKEN_SYMMETRIC_KEY, can run following command)

        ```bash
        openssl rand -hex 64 | head -c 32
        ```
    
    b. create AWS Secrets Manager

    c. move all variables feom .env to AWS Secrets Manager

    d. create step in deploy ci to auto override the local env variables in app.env with prod variables from aws secret manager:

        ```bash
        # download aws cli in you pc

        # run check install  
        aws --version

        # go to aws IAM to generate one access key for the user for this local cli use. 

        # set up credentials to access aws account
        aws configure
        
        AWS Access Key ID [None]: ${{ access_key_id_from_IAM }}
        AWS Secret Access Key [None]: ${{ access_key_password_from_IAM }}
        Default region name [None]: eu-central-1
        Default output format [None]: json
        # your access key id and password will be store under ~/.aws/credentials
        # your region name and output format will be store under ~/.aws/config

        # call secret manager api to retrieve secret values
        aws secretsmanager help
        aws secretsmanager get-secret-value help
        # copy secret arn and secret name from aws secret manager
        # you have to give the user-group of the user in IAM a permission of SecretsManagerReadWrite
        aws secretsmanager get-secret-value --secret-id simple_bank
        # can also use arn: aws secretsmanager get-secret-value --secret-id ${{ arn }} 
        aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text

        # convert json secret value to the format we want and replace local value in app.env with this new variables:
        # use jq package (default package in linux) auto convert json to .env required format:
        aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env
        # result will be:
        # --------
        DB_SOURCE=${{ value from secretsmanager }}
        DB_DRIVER=${{ value from secretsmanager }}
        SERVER_ADDRESS=${{ value from secretsmanager }}
        TOKEN_SYMMETRIC_KEY=${{ value from secretsmanager }}
        ACCESS_TOKEN_DURATION=${{ value from secretsmanager }}
        # --------

        # add the jq command into ci (i.e. deploy.yml)
        - name: Load secrets and save to app.env
          run: aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env

        ```

    e. check the new image in AWS ECR:

       after merge above change without error, new image will be created in ECR. download and run that image:

       ```bash
        # image uri from ecr: 309977797415.dkr.ecr.eu-central-1.amazonaws.com/simplebank:b3ca1521d2187791968ecf667426f0ce9dbeec86

        # login first
        aws ecr get-login-password | docker login --username AWS --password-stdin 309977797415.dkr.ecr.eu-central-1.amazonaws.com  

        # pull image
        docker pull 309977797415.dkr.ecr.eu-central-1.amazonaws.com/simplebank:2a8f83c1991bd3e02f0ca40890732ab2dc3463b1

        # run image
        docker run -p 8080:8080 309977797415.dkr.ecr.eu-central-1.amazonaws.com/simplebank:2a8f83c1991bd3e02f0ca40890732ab2dc3463b1
       ```

4. Kubernetes setup (via AWS EKS):

    a. config cluster in EKS (need to create an EKS-Cluster Role in IAM)

    b. config Node Group (need to create an EC2 Role in IAM - with permission AmazonEKS_CNI_Policy & AmazonEKSWorkerNodePolicy & AmazonEC2ContainerRegistryReadOnly)

    c. install Kubectl:
    ```bash
    ## connect to AWS EKS cluster

    # provide eks permission to User github-ci by creating inline policy of EKS to deployment User Group
    aws eks update-kubeconfig --name simple-bank --region eu-central-1

    # config file will be created under kube folder
    ls -l ~/.kube
    cat ~/.kube/config

    # after checking the context of cluster, connect to this context
    kubectl config use-context arn:aws:eks:eu-central-1:309977797415:cluster/simple-bank

    # check whether can connect to simple-bank cluster
    kubectl cluster-info

    # if get auth error (check https://repost.aws/knowledge-center/amazon-eks-cluster-access):
    # solution: grant access to the cluster to a user who is not the creator of the cluster

    # - check the current user identity (you will find that it is github-ci user, not the root user we used to create the cluster. also access key used is github-ci suer)
    aws sts get-caller-identity
    cat ~/.aws/credentials

    # to My security credentials in AWS, create access key in access keys section. and update aws credentials file:
    vi ~/.aws/credentials
    ####
    [default]
    aws_access_key_id = ${{ new }}
    aws_secret_access_key = ${{ new }}

    [github]
    aws_access_key_id = ${{ old }}
    aws_secret_access_key = ${{ old }}
    ####
    :wq
    
    # check whether still have any auth error:
    kubectl get pods
    kubectl cluster-info

    # now can access the cluster using the credentials of the root user

    # ------ note: tell aws cli to use github credentials
    export AWS_PROFILE=github
    # if windows: $Env:AWS_PROFILE="github"
    # to switch back: export AWS_PROFILE=default

    # But how to allow github user to access the cluster:
    # need to add this user to a special config map
    # 0. you must be in default root profile first: 
    export AWS_PROFILE=default
    # 1. create a folder "eks" in the simple bank project root. And add file "aws-auth.yaml" 
    # 2. apply auth map yaml:
     kubectl apply -f eks/aws-auth.yaml
    # 3. use github credentials
    export AWS_PROFILE=github
    # 4. check:
    kubectl cluster-info


    # install k9s to use kubernetes command easily (compared to kubedctl)
    scoop install k9s
    k9s

    # show all namespaces of the cluster:
    : ns --> enter --> select one namespace to enter and press esc to exit

    # list all service:
    : service --> enter

    list all pods:
    : pods --> enter

    list all cronjobs:
    : cj --> enter

    list all nodes:
    : nodes --> enter

    check config map:
    : configmap --> enter

    to exit k9s:
    : quit --> enter
    ```

5. deploy web app to Kubernetes cluster on AWS EKS (https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)

    a. create deployment.yaml file under folder /eks


    b. run it and check in k9s:
    ```bash
     kubectl apply -f eks/deployment.yaml

     k9s
     :deployment

     # chose first one and click enter to enter pods view. will see pod is not ready. -> press "d" will see failscheduling error in pod
    ```

    c. go to AWS EKS --> Clusters --> simple-bank: will see no nodes in Compute tag: 

        go to Auto Scaling group and change desired capacity from 0 to 1
    

    d. back to pods in k9s, press "d" will see different 2 errors about too many pods:
    
        # go to :nodes --> press "d" --> will see all 4 pods are occupied in Mom-Terminated Pods section (nb of pods depends on VNI you chose when set up the cluster:https://github.com/awslabs/amazon-eks-ami/blob/master/files/eni-max-pods.txt)

        # go to AWS EKS --> clusters --> simple-bank --> Node Group --> will see Instance types is t3.micro

        # Need to delete this one and recreate one with larger instance type (same as step 4)

        # back to k9s --> :deployments --> delete current deployment by ctrl+d

        # rerun:
        kubectl apply -f eks/deployment.yaml

    e. in k9s --> :deployments --> enter (to pods) --> enter (to containers) --> press l and can see log that the server is now listening and serving HTTP requests on port 8080

    f. To enable to send the request to this pod (https://kubernetes.io/docs/concepts/services-networking/service/)

        # create service.yaml under /eks folder: 
        specify type: LoadBalancer; otherwise the default will be ClusterIP and you will have no external IP

        # run:
        kubectl apply -f eks/service.yaml

        # check in k9s
        :services

        # check the external IP (this is actually the domain name of the AWS load balancer service):
        nslookup a9a9507a541a74b40bee0b6162d5ef1c-1729569289.eu-central-1.elb.amazonaws.com

        # using this domain and test api in postman, for exmaple: http://a9a9507a541a74b40bee0b6162d5ef1c-1729569289.eu-central-1.elb.amazonaws.com/users/login


6. Register a domain and set up A-record using AWS Route 53 Dashboard

    a. register a domain from Route 53

    b. route traffics to the Kubernetes cluster:
        click Hosted Zones in Route 53 --> create new record (Address (A) record) --> have domain name like: api.maindomain.com

    c. test routing on this A-record in postman (e.g. http://api.maindomain.com/users/login) 


7. use Ingress to route traffics to different services in Kubernetes:

    Ingress allows us to set up A record only once, but can define multiple rules in the config file to route traffic to different services

    ## a. reset API service to ClusterIP
    change type of the simple-bank API service from LoadBalancer to ClusterIP because we do not want to expose this service to the outside world anymore

    ## b. set up ingress.yaml
    create ingress.yaml in /eks folder

    ## c.  run following commands:
    ```bash
    kubectl apply -f eks/service.yaml
    kubectl apply -f eks/ingress.yaml
    
    #check in k9s --> :ingresses --> no address yet for the ingress
    ```

    ## d. install NGINX Ingress:
    In order for the Ingress resource to work, the cluster must have an ingress controller running.
    
    Unlike other types of controllers which run as part of the kube-controller-manager binary, Ingress controllers are not started automatically with a cluster. 
    (https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)

    (https://kubernetes.github.io/ingress-nginx/deploy/#aws)

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.6.4/deploy/static/provider/aws/deploy.yaml
    ```

    check in k9s --> :pods --> see whether ingress-nginx-controller pod running in the ingress-nginx namespace
    
    check :ingresses --> can see external address to be accessed by the outside world (e.g. ace106ad58eda4d49a07258482038587-a661b702b2f9c978.elb.eu-central-1.amazonaws.com.)

    copy the address to Route 53 hosted zone --> edit original api.maindomain.com A-record

    run following code to check ip addresses of the domain:
    ```bash
    nslookup api.maindomain.com
    nslookup ace106ad58eda4d49a07258482038587-a661b702b2f9c978.elb.eu-central-1.amazonaws.com.
    # both should get same ip addresses
    ```

    in final, test with postman to see whether the the ingress works


8. auto issue TLS certificates with cert-manager and Let's Encrypt

    ## a. download cert-manager:
    (https://cert-manager.io/docs/installation/kubectl/)

    run following code:
    ```bash
    kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml

    # then check in k9s --> :ns --> you should see few pods running for cert-manager namespace
    ```

    ## b. Creating a Basic ACME Issuer:
    (https://cert-manager.io/docs/configuration/acme/)

    create issuer.yaml under /eks folder

    run following command to execute:
    ```bash
    kubectl apply -f eks/issuer.yaml

    # then check in k9s --> :clusterissuers
    # then check in k9s --> :secrets  --> can find its private key
    # then check in k9s --> :certificate  --> still empty
    ```

    ## c. Attached the issuer to ingress
    add following codes to ingress.yaml:
    ```yaml
    annotations:
        cert-manager.io/cluster-issuer: letsencrypt
    

    tls:
    - hosts:
        - api.maindomain.com
        secretName: simple-bank-api-cert
    ```

    then run:
    ```bash
    kubectl apply -f eks/ingress.yaml

    # then check in k9s --> :ns  --> all
    # then check in k9s --> :ingress  --> press "d" --> will see TLS is enabled
    # then check in k9s --> :certificate 
    ```

    test on postman, both should work:
    https://api.maindomain.com/users/login
    http://api.maindomain.com/users/login


9. Auto deploy to kubernetes with github action:

    ## a. update deploy.yml
    (https://github.com/marketplace/actions/kubectl-tool-installer)
    (https://storage.googleapis.com/kubernetes-release/release/stable.txt)

    - add new step: name: Install kubectl
    - update steo: name: Build, tag, and push docker image to Amazon ECR
        ```bash
        # old
        run: |
            docker build -t $REGISTRY/$REPOSITORY:$IMAGE_TAG .
            docker push $REGISTRY/$REPOSITORY:$IMAGE_TAG
        
        # new
        run: |
            docker build -t $REGISTRY/$REPOSITORY:$IMAGE_TAG -t latest .
            docker push -a $REGISTRY/$REPOSITORY:$IMAGE_TAG 
        ```
    
    - add new step: name: Deploy image to Amazon EKS

    ## b. update deployment.yaml
    ```bash
    # old
    spec:
      containers:
      - name: simple-bank-api
        image: 309977797415.dkr.ecr.eu-central-1.amazonaws.com/simplebank:2a8f83c1991bd3e02f0ca40890732ab2dc3463b1
        ports:
        - containerPort: 8080
    
    # new
    spec:
      containers:
      - name: simple-bank-api
        # using latest instead of image tag
        image: 309977797415.dkr.ecr.eu-central-1.amazonaws.com/simplebank:latest
        ports:
        - containerPort: 8080
    ```