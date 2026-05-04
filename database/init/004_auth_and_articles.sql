create table if not exists user_service.credentials (
    user_id text primary key references user_service.users(id) on delete cascade,
    email text not null unique,
    password_hash text not null,
    created_at timestamptz not null default now()
);

create table if not exists gamification.articles (
    id text primary key,
    author_id text not null references user_service.users(id) on delete cascade,
    title text not null,
    summary text not null,
    body text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists gamification.article_comments (
    id text primary key,
    article_id text not null references gamification.articles(id) on delete cascade,
    author_id text not null references user_service.users(id) on delete cascade,
    body text not null,
    created_at timestamptz not null default now()
);

create table if not exists gamification.article_reactions (
    article_id text not null references gamification.articles(id) on delete cascade,
    user_id text not null references user_service.users(id) on delete cascade,
    reaction_type text not null,
    created_at timestamptz not null default now(),
    primary key (article_id, user_id, reaction_type)
);

create index if not exists idx_articles_created_at on gamification.articles(created_at desc);
create index if not exists idx_article_comments_article_id on gamification.article_comments(article_id);
create index if not exists idx_article_reactions_article_id on gamification.article_reactions(article_id);
