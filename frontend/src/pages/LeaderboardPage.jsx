import { useAppData } from "../context/AppDataContext";

const podiumOrder = [1, 0, 2];

function LeaderboardPage() {
  const { portalState } = useAppData();
  const { podium, rows } = portalState.leaderboard;

  return (
    <section className="page">
      <div className="page__intro">
        <div>
          <h1>Рейтинг разработчиков</h1>
          <p>Соревнуйтесь и поднимайтесь в рейтинге.</p>
        </div>
      </div>

      <section className="surface-card">
        <div className="podium">
          {podiumOrder.map((index) => {
            const item = podium[index];
            if (!item) {
              return null;
            }

            const height = item.rank === 1 ? 102 : item.rank === 2 ? 78 : 64;

            return (
              <article key={item.userId} className={`podium__item podium__item--rank-${item.rank}`}>
                <div className="podium__avatar">{item.name.slice(0, 1)}</div>
                <strong>{item.name}</strong>
                <span>{item.title}</span>
                <em>{item.levelText}</em>
                <div className="podium__bar" style={{ height: `${height}px` }} />
                <p>{new Intl.NumberFormat("ru-RU").format(item.xp)} XP</p>
              </article>
            );
          })}
        </div>
      </section>

      <section className="surface-card">
        <div className="table-wrap">
          <table className="leaderboard-table">
            <thead>
              <tr>
                <th>Ранг</th>
                <th>Разработчик</th>
                <th>Компания</th>
                <th>Уровень</th>
                <th>XP</th>
              </tr>
            </thead>

            <tbody>
              {rows.map((row) => (
                <tr key={row.userId} className={row.isCurrent ? "leaderboard-table__current" : ""}>
                  <td>{row.rank <= 3 ? `Топ ${row.rank}` : `#${row.rank}`}</td>
                  <td>
                    <div className="leaderboard-user">
                      <span className="leaderboard-user__dot">{row.name.slice(0, 1)}</span>
                      <div>
                        <strong>{row.name}</strong>
                        <p>{row.title}</p>
                      </div>
                    </div>
                  </td>
                  <td>{row.company}</td>
                  <td>{row.levelText}</td>
                  <td>{new Intl.NumberFormat("ru-RU").format(row.xp)} XP</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </section>
  );
}

export default LeaderboardPage;
