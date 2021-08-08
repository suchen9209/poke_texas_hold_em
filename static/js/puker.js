// 扑克牌
puker = {
    AllCards : [
        'puker-spade14','puker-spade2','puker-spade3','puker-spade4','puker-spade5','puker-spade6','puker-spade7','puker-spade8','puker-spade9','puker-spade10','puker-spade11','puker-spade12','puker-spade13',
        'puker-heart14','puker-heart2','puker-heart3','puker-heart4','puker-heart5','puker-heart6','puker-heart7','puker-heart8','puker-heart9','puker-heart10','puker-heart11','puker-heart12','puker-heart13',
        'puker-diamond14','puker-diamond2','puker-diamond3','puker-diamond4','puker-diamond5','puker-diamond6','puker-diamond7','puker-diamond8','puker-diamond9','puker-diamond10','puker-diamond11','puker-diamond12','puker-diamond13',
        'puker-club14','puker-club2','puker-club3','puker-club4','puker-club5','puker-club6','puker-club7','puker-club8','puker-club9','puker-club10','puker-club11','puker-club12','puker-club13'
        ],
    KingsB : ['puker-big-kingB','puker-small-kingB'],
    KingsA : ['puker-big-kingA','puker-small-kingA'],  
}
//随机选择n张牌
puker.chooseCards = function(n){
    puker.shuffle(puker.AllCards);
    return puker.AllCards.slice(0,n);
}

//洗牌
puker.shuffle = function(cards){
    cards.sort(function(){
        return 0.5-Math.random();
    });
    return cards;
};