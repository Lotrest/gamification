import '../CSS/ProfilePagesCSS.css'

function ProfilePages() {
    return (
        <div className='profile-area'>

            <div className='dashbord-top'>

                <div className='dashbord-left'>
                    <h1>Профиль</h1>
                    <p>Сделать вклад в комьюнити не только полезным, но и захватывающим</p>
                </div>

                <div className='dashbord-right'>
                    <div>
                        <img width="40" height="40" src="https://img.icons8.com/ios/50/cheap-2.png" alt="cheap-2" />
                    </div>

                    <div>
                        <img width="40" height="40" src="https://img.icons8.com/external-linear-outline-icons-papa-vector/78/external-Notifications-interface-linear-outline-icons-papa-vector.png" alt="vot" />
                    </div>

                    <div className='mini-user'>
                        <img width="40" height="40" src="src/img/user.png" alt="user" />
                        <div>
                            <p>Пользователь</p>
                            <p>Уровень</p>
                        </div>
                    </div>
                </div>

            </div>

            <div className='size'>


                <div className='conteiner-user'>

                    <div className='conteiner-user-left'>
                        <img width="70" height="70" src="src/img/user.png" alt="user" />
                        <div>
                            <h3>Пользователь</h3>
                            <p className='color'>Должность</p>
                            <p className='color'>На портале с Марта 2023 &nbsp; | &nbsp; Ранг #10 из 1 248</p>
                        </div>
                    </div>

                    <div className='conteiner-user-right'>
                        <div className='circle'>
                            <p className='lv'>1ур.</p>
                            <p className='exp'>40%</p>
                        </div>
                        <p>200XP</p>
                    </div>
                </div>

                <div className='conteiner-stats'>
                    
                    <div className='conteiner-stats-fill'>
                        <p>Api запросы</p>
                    </div>

                    <div className='conteiner-stats-fill'>
                        <p>Статьи</p>
                    </div>

                    <div className='conteiner-stats-fill'>
                        <p>Комментарии</p>
                    </div>

                    <div className='conteiner-stats-fill'>
                        <p>Всего xp</p>
                    </div>

                </div>

                <div className='articles-and-achievements'>

                    <div className='conteiner-articles'>
                        <div className='heading-button'>
                            <h3>Топ статей</h3>
                            <button>Все статьи</button>
                        </div>
                        <div className='public-article'>
                            <div>
                                <p className='public-article-p-1'>Интеграция с API CDEK: Полный гайд</p>
                                <p className='public-article-p-2'>2 340 просмотров • 18 комментариев • +200 XP</p>
                            </div>
                            <div className='star'>
                                <img width="30" height="30" src="https://img.icons8.com/arcade/64/star.png" alt="star"/>
                                <p>4.9</p>
                            </div>
                        </div>
                    </div>

                    <div className='conteiner-achievements'>
                        <h3>Бэйджи</h3>

                        <div className='achievements'>
                            <div className='achievements-stile'>
                                <p className='public-article-p-1'>Заркий глаз</p>
                                <p className='public-article-p-2'>Обычный</p>
                            </div>

                            <div className='achievements-stile'>
                                <p className='public-article-p-1'>Эмпат</p>
                                <p className='public-article-p-2'>Обычный</p>
                            </div>

                            <div className='achievements-stile'>
                                <p className='public-article-p-1'>Слухач</p>
                                <p className='public-article-p-2'>Обычный</p>
                            </div>

                        </div>
                    </div>

                </div>

                <div className='conteiner-recent-actions'>
                    <div>
                        <h3>Последние действия</h3>
                    </div>
                    <div className='recent-actions'>
                        <div className='recent-actions-option'>
                            <img width="48" height="48" src="https://img.icons8.com/emoji/48/check-mark-emoji.png" alt="check-mark-emoji"/>
                            <div>
                                <p className='actions'>Получен бэйдж "Детектив логов"</p>
                                <p className='time'>2 часа назад</p>
                            </div>
                        </div>
                        <div className='xp'>   
                            <p >+90 XP</p>
                        </div>
                    </div>
                </div>


            </div>


        </div>
    )
}
export default ProfilePages;