package database

const schema = `
CREATE TABLE IF NOT EXISTS cliente_status (
    id SERIAL PRIMARY KEY,
    descricao VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS cliente_prioridades (
    id SERIAL PRIMARY KEY,
    descricao VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS clientes (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    tipo_solicitacao VARCHAR(100) NOT NULL,
    valor_patrimonio BIGINT NOT NULL,
    status_id INT NOT NULL DEFAULT 1,
    prioridade_id INT DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_cliente_status
        FOREIGN KEY (status_id) REFERENCES cliente_status(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_cliente_prioridade
        FOREIGN KEY (prioridade_id) REFERENCES cliente_prioridades(id)
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_clientes_email ON clientes(email);

CREATE TABLE IF NOT EXISTS eventos_processados (
    event_id VARCHAR(255) PRIMARY KEY,
    card_id VARCHAR(255) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO cliente_status (id, descricao) VALUES
(1, 'Aguardando Análise'),
(2, 'Processado')
ON CONFLICT (id) DO NOTHING;

INSERT INTO cliente_prioridades (id, descricao) VALUES
(1, 'prioridade_normal'),
(2, 'prioridade_alta')
ON CONFLICT (id) DO NOTHING;
`
