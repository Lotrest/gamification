create table if not exists user_service.users (
    id text primary key,
    name text not null,
    title text not null,
    company text not null,
    level int not null default 1,
    level_text text not null default 'Lv.1',
    joined_at text not null,
    location text not null,
    team text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists gamification.user_state (
    user_id text primary key,
    current_xp int not null default 0,
    coins int not null default 0,
    today_earned int not null default 0,
    streak_days int not null default 0,
    rank int not null default 1,
    completed_tasks int not null default 0,
    api_requests int not null default 0,
    articles_count int not null default 0,
    comments_count int not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists gamification.weekly_activity (
    user_id text not null,
    day_code text not null,
    xp int not null default 0,
    sort_order int not null,
    primary key (user_id, day_code)
);

create table if not exists gamification.weekly_activity (
    user_id text not null,
    day_code text not null,
    xp int not null default 0,
    sort_order int not null,
    primary key (user_id, day_code)
);

create table if not exists gamification.recent_activity (
    id text primary key,
    user_id text not null,
    title text not null,
    timestamp_label text not null,
    xp int not null default 0,
    sort_order int not null default 0
);

create table if not exists gamification.tasks (
    id text primary key,
    user_id text not null,
    title text not null,
    description text not null,
    status text not null,
    progress int not null default 0,
    target int not null default 1,
    reward_xp int not null default 0
);

create table if not exists gamification.achievements (
    id text primary key,
    user_id text not null,
    title text not null,
    description text not null,
    rarity text not null,
    status text not null,
    reward_xp int not null default 0
);

create table if not exists gamification.leaderboard (
    user_id text primary key,
    rank int not null,
    xp int not null default 0
);

create table if not exists gamification.rewards (
    id text primary key,
    user_id text not null,
    title text not null,
    description text not null,
    cost int not null,
    status text not null,
    category text not null
);

create table if not exists gamification.purchases (
    id text primary key,
    user_id text not null,
    reward_id text not null,
    title text not null,
    cost int not null,
    redeemed_at text not null,
    status text not null
);

create table if not exists gamification.notifications (
    id text primary key,
    user_id text not null,
    title text not null,
    body text not null,
    variant text not null,
    created_at timestamptz not null default now()
);
