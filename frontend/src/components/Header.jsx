import '../CSS/HeaderCSS.css';
import { NavLink } from "react-router-dom";

function Header(){
    return(
        <section className="">
            <div className='header-img'>
                <img src="src/img/logo.png" alt="logo" />
            </div>
            <nav className='header-nav'>
                <p>Навигация</p>
                <NavLink to="/">Профиль</NavLink>
                <NavLink to="/dashbord">Дашборд</NavLink>
                <NavLink to="/achievements">Достижения</NavLink>
                <NavLink to="/tasks">Задания</NavLink>
                <NavLink to="/runked">Рейтинг</NavLink>
                <NavLink to="/shop">Магазин наград</NavLink>
            </nav>
        </section>
    )
}
export default Header