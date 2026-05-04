import '../CSS/DashbordPagesCSS.css'

function DashbordPages(){
    return(
        <div className='profile-area'>
            <div className='dashbord-top'>

                <div className='dashbord-left'>
                    <h1>Дашборд</h1>
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

                <div className='conteiner-dashbord'>

                    <div className='dashbord-left-stile'>

                        <img width="40" height="40" src="src/img/user.png" alt="user" />

                        <div className='test'>
                            <p className='option-text-1'>Пользователь</p>
                            <p className='option-text-2'>Должность</p>
                            <p className='option-text-2'>уровень</p>
                            <div className='progress'>
                                <div className='progress-bar'>

                                </div>
                            </div>
                            <p className='option-text-2'>До уровня</p>
                        </div>
                    </div>

                    <div className='dashbord-right-stile'>
                        <div className='dashbord-right-stile-color'>
                            <p className='option-text-1'>#12</p>
                            <p className='option-text-2'>ранг из 1248</p>
                        </div>
                        <div className='dashbord-right-stile-color'>
                            <p className='option-text-1'>14</p>
                            <p className='option-text-2'>Бейджей</p>
                        </div>
                        <div className='dashbord-right-stile-color'>
                            <p className='option-text-1'>7 д.</p>
                            <p className='option-text-2'>Серия подряд</p>
                        </div>
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


                <div className='conteiner-activity-and-task'>
                    <div className='conteiner-activity'>
                        <div className='conteiner-h3-button'>
                            <h3>Активность за неделю</h3>
                            <button>Полную статистику</button>
                        </div>
                        
                        <div className='conteiner-xp-day'>
                            <div className='bar-item'>
                                <span>XP</span>
                                <div  className="bar" style={{ height: "60px" }}></div>
                                <p>День</p>
                            </div>
                            <div className='bar-item'>
                                <span>XP</span>
                                <div  className="bar" style={{ height: "70px" }}></div>
                                <p>День</p>
                            </div>
                            <div className='bar-item'>
                                <span>XP</span>
                                <div  className="bar" style={{ height: "50px" }}></div>
                                <p>День</p>
                            </div>
                            <div className='bar-item'>
                                <span>XP</span>
                                <div  className="bar" style={{ height: "40px" }}></div>
                                <p>День</p>
                            </div>
                            <div className='bar-item'>
                                <span>XP</span>
                                <div  className="bar" style={{ height: "20px" }}></div>
                                <p>День</p>
                            </div>
                            <div className='bar-item'>
                                <span>XP</span>
                                <div  className="bar" style={{ height: "30px" }} ></div>
                                <p>День</p>
                            </div>
                            <div className='bar-item'>
                                <span>XP</span>
                                <div  className="bar" style={{ height: "90px" }}></div>
                                <p>День</p>
                            </div> 
                        </div>

                        <hr />

                         <div className="chart-footer">
                            <p>Итог за неделю: <span>1780 XP</span></p>
                            <p>Лучший день: <span>Пт — 250 XP</span></p>
                        </div>
                    </div>

                    <div className='conteiner-task'>
                        <p>Активные заданий</p>
                    </div>
                </div>

            </div>

        </div>
    )
}
export default DashbordPages;