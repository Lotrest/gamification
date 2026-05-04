insert into user_service.users (id, name, title, company, level, level_text, joined_at, location, team)
values ('me', 'Пользователь', 'QA-инженер', 'CDEK Digital', 1, 'Lv.1', 'сегодня', 'Новосибирск', 'Платформа комьюнити')
on conflict (id) do nothing;

insert into gamification.user_state (
    user_id,
    current_xp,
    coins,
    today_earned,
    streak_days,
    rank,
    completed_tasks,
    api_requests,
    articles_count,
    comments_count
)
values ('me', 0, 0, 0, 0, 1, 0, 0, 0, 0)
on conflict (user_id) do nothing;

insert into gamification.weekly_activity (user_id, day_code, xp, sort_order)
values
    ('me', 'Пн', 0, 1),
    ('me', 'Вт', 0, 2),
    ('me', 'Ср', 0, 3),
    ('me', 'Чт', 0, 4),
    ('me', 'Пт', 0, 5),
    ('me', 'Сб', 0, 6),
    ('me', 'Вс', 0, 7)
on conflict (user_id, day_code) do nothing;

insert into gamification.tasks (id, user_id, title, description, status, progress, target, reward_xp)
values
    ('task-first-article', 'me', 'Задание на статью', 'Напиши и опубликуй первую статью в базе знаний платформы.', 'available', 0, 1, 150),
    ('task-scouting', 'me', 'Скаутинг', 'Сделай 3 полезных действия под статьями: лайк, репост или комментарий.', 'available', 0, 3, 80),
    ('task-feedback', 'me', 'Корректор', 'Оставь один содержательный комментарий к статье.', 'available', 0, 1, 40)
on conflict (id) do nothing;

insert into gamification.achievements (id, user_id, title, description, rarity, status, reward_xp)
values
    ('achievement-eye', 'me', 'Зоркий глаз', 'Сделай первое полезное действие на платформе.', 'Обычные', 'locked', 20),
    ('achievement-empathy', 'me', 'Эмпат', 'Оставь первый комментарий к статье.', 'Обычные', 'locked', 10),
    ('achievement-author', 'me', 'Автор', 'Опубликуй первую статью.', 'Редкие', 'locked', 70),
    ('achievement-guardian', 'me', 'Страж', 'Сделай 3 полезных действия под статьями.', 'Эпические', 'locked', 150)
on conflict (id) do nothing;

insert into gamification.leaderboard (user_id, rank, xp)
values ('me', 1, 0)
on conflict (user_id) do nothing;

insert into gamification.rewards (id, user_id, title, description, cost, status, category)
values
    ('reward-hoodie', 'me', 'Фирменный худи', 'Мерч команды платформы.', 2500, 'available', 'мерч'),
    ('reward-coffee', 'me', 'Кофе с архитектором', 'Часовая встреча и разбор архитектурных решений.', 1800, 'available', 'бонус')
on conflict (id) do nothing;

insert into gamification.notifications (id, user_id, title, body, variant)
values (
    'notification-welcome',
    'me',
    'Добро пожаловать!',
    'Платформа запущена. Пока прогресс нулевой — начни с первого задания.',
    'success'
)
on conflict (id) do nothing;
