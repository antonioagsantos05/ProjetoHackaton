-- Inicialização do Banco de Dados PostgreSQL FIAP X
-- Baseado no Documento Técnico de Banco de Dados FIAP X

-- Tabela de Usuário
CREATE TABLE usuario (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    email TEXT NOT NULL,
    nome TEXT NOT NULL,
    hash_senha TEXT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE (tenant_id, email)
);
CREATE INDEX idx_usuario_status ON usuario USING btree (status);
CREATE INDEX idx_usuario_created_at ON usuario (created_at);

-- Tabela de Sessão
CREATE TABLE sessao (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    usuario_id UUID NOT NULL REFERENCES usuario(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_sessao_tenant_usuario ON sessao (tenant_id, usuario_id);
CREATE INDEX idx_sessao_expires_at ON sessao (expires_at);

-- Tabela de Vídeo
CREATE TABLE video (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    usuario_id UUID REFERENCES usuario(id) ON DELETE SET NULL,
    titulo TEXT NOT NULL,
    video_path TEXT, -- Caminho do objeto no storage (Minio)
    status SMALLINT NOT NULL DEFAULT 0,
    meta JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_video_tenant_status_ct ON video (tenant_id, status, created_at);
CREATE INDEX idx_video_meta_gin ON video USING gin (meta);

-- Tabela de Worker
CREATE TABLE worker (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    nome TEXT UNIQUE NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    last_heartbeat TIMESTAMPTZ,
    capacidades JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_worker_status ON worker (status);
CREATE INDEX idx_worker_capacidades_gin ON worker USING gin (capacidades);

-- Tabela de Job Processamento
CREATE TABLE job_processamento (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    video_id UUID NOT NULL REFERENCES video(id) ON DELETE CASCADE,
    worker_id UUID REFERENCES worker(id) ON DELETE SET NULL,
    tipo SMALLINT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    params JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_job_tenant_status_started ON job_processamento (tenant_id, status, started_at);
CREATE INDEX idx_job_tenant_video_created ON job_processamento (tenant_id, video_id, created_at);
CREATE INDEX idx_job_params_gin ON job_processamento USING gin (params);

-- Tabela de Frame
CREATE TABLE frame (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    video_id UUID NOT NULL REFERENCES video(id) ON DELETE CASCADE,
    frame_no BIGINT NOT NULL,
    timestamp_ms BIGINT NOT NULL,
    features JSONB,
    uri TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE (tenant_id, video_id, frame_no)
);
CREATE INDEX idx_frame_timestamp ON frame (timestamp_ms);
CREATE INDEX idx_frame_features_gin ON frame USING gin (features);

-- Tabela de Arquivo ZIP
CREATE TABLE arquivo_zip (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    job_id UUID NOT NULL REFERENCES job_processamento(id) ON DELETE CASCADE,
    uri TEXT NOT NULL,
    checksum TEXT,
    tamanho BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_zip_tenant_job ON arquivo_zip (tenant_id, job_id);

-- Tabela de Notificação
CREATE TABLE notificacao (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    usuario_id UUID NOT NULL REFERENCES usuario(id) ON DELETE CASCADE,
    tipo SMALLINT NOT NULL,
    payload JSONB,
    lida BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_notificacao_tenant_usuario_lida ON notificacao (tenant_id, usuario_id, lida);
CREATE INDEX idx_notificacao_payload_gin ON notificacao USING gin (payload);

-- Tabela de Política Retenção
CREATE TABLE politica_retencao (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    scope SMALLINT NOT NULL,
    dias INTEGER NOT NULL CHECK (dias > 0),
    ativo BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_politica_tenant_scope_ativo ON politica_retencao (tenant_id, scope, ativo) WHERE ativo;

-- Exemplo RLS (Row Level Security) - Deve ser ativado por tabela
ALTER TABLE video ENABLE ROW LEVEL SECURITY;
ALTER TABLE job_processamento ENABLE ROW LEVEL SECURITY;
-- e assim por diante para todas as tabelas...

-- Politica basica RLS
-- CREATE POLICY tenant_isolation_policy ON video
--     USING (tenant_id = current_setting('app.tenant_id')::uuid);
