import { useState, useEffect } from "react";
import { Elements, PaymentElement, useStripe, useElements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';
import Router from 'next/router';

// Make sure to call `loadStripe` outside of a component’s render to avoid
// recreating the `Stripe` object on every render.
const stripePromise = loadStripe('pk_test_51PniCxFDYBNcFTueZtWj2bvZ0OCfWKXXHN6FrE2tjGJZuctK6w2b1O52SUyhKFuuIXH4FtmOvoDexljkCCqAjBmI00T4jT2TCX');

const Order = ({ order, currentUser }) => {

    const LeftTime = () => {
        const calTimeLeft = () => new Date(order.expiresAt) - new Date()
        const [msLeft, setMsLeft] = useState(calTimeLeft());

        useEffect(() => {
            const timer = setInterval(() => setMsLeft(calTimeLeft()), 1000);
            return () => { clearInterval(timer) }
        }, [order]);
        if (msLeft <= 0) {
            return <p>Order Expired</p>;
        }
        return <p>Time left to pay: {msLeft} seconds</p>;
    };

    const SetupForm = () => {
        const stripe = useStripe();
        const elements = useElements();

        const [errorMessage, setErrorMessage] = useState();
        const [loading, setLoading] = useState(false);

        const handleError = (error) => {
            setLoading(false);
            setErrorMessage(error.message);
        }
        const handleSubmit = async (event) => {
            // We don't want to let default form submission happen here,
            // which would refresh the page.
            event.preventDefault();

            if (!stripe) {
                // Stripe.js hasn't yet loaded.
                // Make sure to disable form submission until Stripe.js has loaded.
                return;
            }

            setLoading(true);

            // Trigger form validation and wallet collection
            const { error: submitError } = await elements.submit();
            if (submitError) {
                handleError(submitError);
                return;
            }

            // Create the PaymentIntent and obtain clientSecret
            const res = await fetch("/api/payments", {
                method: "POST",
                body: JSON.stringify({ orderId: order.orderID })
            });

            const { clientSecret } = await res.json();

            // Confirm the PaymentIntent using the details collected by the Payment Element
            const { error } = await stripe.confirmPayment({
                elements,
                clientSecret,
                confirmParams: {
                    return_url: 'https://cca8-2405-9800-b950-4e0a-a829-2af5-f7cb-a9da.ngrok-free.app/orders',
                },
            });

            if (error) {
                // This point is only reached if there's an immediate error when
                // confirming the setup. Show the error to your customer (for example, payment details incomplete)
                handleError(error);
            } else {
                // Your customer is redirected to your `return_url`. For some payment
                // methods like iDEAL, your customer is redirected to an intermediate
                // site first to authorize the payment, then redirected to the `return_url`.
            }
        };
        return (
            <form onSubmit={handleSubmit}>
                <PaymentElement />
                <button type="submit" disabled={!stripe || loading}>
                    Submit
                </button>
                {errorMessage && <div>{errorMessage}</div>}
            </form>
        );
    };

    const options = {
        mode: 'payment',
        amount: order.ticket.price,
        currency: 'usd',
        // Fully customizable with appearance API.
        appearance: {/*...*/ },
    };

    return (
        <div>
            <LeftTime />
            <Elements stripe={stripePromise} options={options}>
                <SetupForm />
            </Elements>
        </div>
    );
}


Order.getInitialProps = async (context, client, currentUser) => {
    const { data } = await client.get(`/api/orders/${context.query.id}`);
    return { order: data };
}

export default Order;