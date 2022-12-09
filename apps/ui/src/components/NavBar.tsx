import React from "react";
import { Link } from "react-router-dom";

const NavBar: React.FunctionComponent<React.PropsWithChildren<unknown>> = () => {
	return (
		<div>
			<nav>
				<ul>
					<li>
						<Link to="/">Timeline</Link>
					</li>
					<li>
						<Link to="/">Places</Link>
					</li>
					<li>
						<Link to="/">People</Link>
					</li>
					<li>
						<Link to="/">Albums</Link>
					</li>
				</ul>
			</nav>
			{/* <div>Hello, {user.username}</div>
                <AmplifySignOut /> */}
		</div>
	);
};

export default NavBar;
