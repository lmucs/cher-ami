## User
  -
  - Info: pas...
  - :WROTE                  -> Message
  - :CHIEF_OF               -> Circle

## Message
  - Content
  - :PUB_TO                 -> Circle

## Circles
  - :INCLUDES               -> Users


CREATE (Brovek  :User { data: {...} })
CREATE (Salamat :User { data: {...} })
CREATE (Otto    :User { data: {...} })

CREATE (msg_1:Messages { content : 'Water is the best drink.' })
CREATE (msg_2:Messages { content : 'Watermelon is not the best fruit.' })

CREATE (Bro_gold : Circle {})
CREATE (Bro_pub  : Circle {})
CREATE (Sal_gold : Circle {})

--- connections

CREATE (Brovek)   -[:WROTE {}]       -> (msg_1)
CREATE (Salamat)  -[:WROTE {}]       -> (msg_2)

CREATE (Brovek)   -[:CHIEF_OF {}]    -> (Bro_gold)
CREATE (Brovek)   -[:CHIEF_OF {}]    -> (Bro_pub)
CREATE (Salamat)  -[:CHIEF_OF {}]    -> (Sal_gold)

CREATE (msg_1)      -[:PUB_TO {}]   -> (Bro_pub)
CREATE (msg_2)      -[:PUB_TO {}]   -> (Sal_gold)

CREATE (Bro_gold)   -[:INCLUDES {}]    -> (Salamat)
CREATE (Bro_pub)    -[:INCLUDES {}]    -> (Salamat)
CREATE (Bro_pub)    -[:INCLUDES {}]    -> (Otto)
CREATE (Sal_gold)   -[:INCLUDES {}]    -> (Otto)
