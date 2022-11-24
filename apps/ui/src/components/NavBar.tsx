import React from 'react';
import { Link } from 'react-router-dom';

const NavBar: React.FunctionComponent = () => {
    return (

        < div >
            <nav>
                <Link to="/">Timeline</Link>
                {/* <div>Hello, {user.username}</div>
                <AmplifySignOut /> */}
            </nav>
        </div >

    )
}

export default NavBar;
