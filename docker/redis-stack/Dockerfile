FROM redis/redis-stack:latest

# Installer les dépendances de construction
RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    pkg-config \
    libssl-dev \
    clang \
    libclang-dev

# Installer Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

# Configurer l'environnement Rust
ENV PATH="/root/.cargo/bin:${PATH}"

# Cloner le dépôt RedisJSON
RUN git clone https://github.com/RedisJSON/RedisJSON.git /usr/src/redisjson

# Construire RedisJSON
RUN cd /usr/src/redisjson && \
    cargo build --release

# Configurer Redis pour charger le module RedisJSON
RUN echo "loadmodule /usr/src/redisjson/target/release/librejson.so" >> /redis-stack.conf

# Configurer le mot de passe d'authentification pour Redis
RUN echo "requirepass training" >> /redis-stack.conf

# Exposer le port Redis
EXPOSE 6379

# Démarrer Redis avec le module RedisJSON
CMD ["redis-server", "/redis-stack.conf"]
