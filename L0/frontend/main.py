
import streamlit as st
import requests

API_BASE = "http://order_service:8080"

st.title("Order Info")

order_uid = st.text_input("Введите order_uid:")

if st.button("Получить заказ"):
    if not order_uid:
        st.warning("Введите order_uid")
    else:
        url = f"{API_BASE}/orders/{order_uid}"
        try:
            res = requests.get(url)
            if res.status_code == 200:
                data = res.json()
                st.success("Заказ найден!")
                st.json(data)  # красиво подсвеченный JSON
            elif res.status_code == 404:
                st.error("Заказ не найден (404)")
            else:
                st.error(f"Ошибка: {res.status_code}")
        except Exception as e:
            st.error(f"Ошибка запроса: {e}")
