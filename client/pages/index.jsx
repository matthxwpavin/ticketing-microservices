import Link from 'next/link';

const LandingPage = ({ currentUser, tickets }) => {
  const ticketList = tickets.map(ticket => (
    <tr key={ticket.id}>
      <td>{ticket.title}</td>
      <td>{ticket.price / 100}</td>
      <td>
        <Link href={`/tickets/${ticket.id}`}>View</Link>
      </td>
    </tr>
  ));
  return (
    <div>
      <h1>Tickets</h1>
      <table className="table">
        <thead>
          <tr>
            <th>Title</th>
            <th>Price</th>
            <th>Link</th>
          </tr>
        </thead>
        <tbody>
          {ticketList.reverse()}
        </tbody>
      </table>
    </div>
  );
};

LandingPage.getInitialProps = async (context, client, currentUser) => {
  let tickets = [];
  try {
    const { data } = await client.get('/api/tickets').catch();
    tickets = data;
  } catch { }
  return { tickets };
};

export default LandingPage;
