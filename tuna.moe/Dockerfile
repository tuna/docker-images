FROM tunathu/mirror-web
  
RUN apt-get install --no-install-recommends -y npm
RUN git clone https://github.com/developmentseed/jekyll-hook.git /jekyll-hook
RUN cd /jekyll-hook && npm install
RUN cd /jekyll-hook && git revert -n  0de7bcd27a129d27d21af88e6753e4915a626f8a && echo '/Td6WFoAAATm1rRGAgAhARYAAAB0L+Wj4AJcAU9dADIaSQnC/BF9UN4KT0fXfE0JwDsnPrWRkkrP0gDAqKBEk3qhI9nn6M9T8DsdwJCVQPL8b8vFMfyB/jdbTDaBjSppL3nAG/sACpu0NOae93qTpQeeTr6CoLUMWffeD8jZ1E34JJTgW8SNDezdgkbsdDa2x9iv8S7Zih0vzqItL64kC5+f2xOKOS8InT+QzF2+bj/Woy/LRztY9AUl30DfqHwbH742M4f1AXHjWW0hLPoGbMaC4unwGXix0IkfCd5+Oox7dBmk/GmrW6WcIg0euO0GXcldLbj+0GrQD0OK5g1RXImZmwS8/msXQv5aGTf99ffNBCUE20ieVJE9QXFoAE/3llDNzrbrKh2P1h16Ti3TvQqLr8oi6ha05DHb6dKVcIpHwdD8bkHCDbOmEf5HnZT42SNnqxQoo/zwHQJO2pzctbp9w8QbpY6OKAbR542TAAC3+IFSSOChKAAB6wLdBAAAZKUhFrHEZ/sCAAAAAARZWg==' | base64 -d | xz -d | patch

EXPOSE 8080
WORKDIR /jekyll-hook

CMD ["nodejs", "/jekyll-hook/jekyll-hook.js"]
