FROM redis/redis-stack:latest

# Installer les dépendances de construction
RUN apt-get update && apt upgrade -y && apt-get install -y \
    curl \
    wget \
    unzip \
    git \
    build-essential \
    pkg-config \
    libssl-dev \
    libtool \
    automake \
    autoconf \
    clang \
    libclang-dev \
    cmake \
    m4 \
    peg \
    python3 \
    python3-pip

# Installer Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

# Configurer l'environnement Rust
ENV PATH="/root/.cargo/bin:${PATH}"

# Cloner les dépôts des modules Redis
RUN git clone https://github.com/RedisJSON/RedisJSON.git /usr/src/redisjson
RUN git clone --branch v2.6.3 https://github.com/RediSearch/RediSearch.git /usr/src/redisearch
RUN git clone --recursive https://github.com/RedisAI/RedisAI.git /usr/src/redisai
RUN git clone --recurse-submodules -j8 https://github.com/RedisGraph/RedisGraph.git /usr/src/redisgraph



# Construire les modules Redis
RUN cd /usr/src/redisjson && cargo build --release
RUN cd /usr/src/redisearch && make build
RUN cd /usr/src/redisai && bash get_deps.sh && make -C opt ALL=1 && make -C opt
RUN cd /usr/src/redisgraph && make

# Configurer Redis pour charger les modules
RUN echo "loadmodule /usr/src/redisjson/target/release/librejson.so" >> /redis-stack.conf
RUN echo "loadmodule /usr/src/redisearch/bin/linux-x64-release/search/redisearch.so" >> /redis-stack.conf
RUN echo "loadmodule /usr/src/redisai/bin/linux-x64-release/src/redisai.so" >> /redis-stack.conf
RUN echo "loadmodule /usr/src/redisgraph/bin/linux-x64-release/src/redisgraph.so" >> /redis-stack.conf

# Configurer le mot de passe d'authentification pour Redis
RUN echo "requirepass training" >> /redis-stack.conf

# Exposer le port Redis
EXPOSE 6379

# Démarrer Redis avec le module RedisJSON
CMD ["redis-server", "/redis-stack.conf"]
