#!/bin/bash
NAMES=(Ana Bruno Carla Diego Elena Felipe Gabriela Hugo Isabela Joao Karen Lucas Mariana Nicolas Olivia Pedro Rafael Sofia Thiago Camila)
SURNAMES=(Silva Santos Oliveira Souza Lima Pereira Costa Ferreira Alves Rodrigues Martins Gomes Barbosa Carvalho Rocha Dias Nascimento Andrade Moreira Nunes)
POSITIONS=(Engineer Designer Manager Analyst DevOps QA Architect)
TYPES=(fulltime contractor)

echo "Inserindo 400 registros..."
for i in $(seq 1 400); do
  NAME="${NAMES[$RANDOM % ${#NAMES[@]}]} ${SURNAMES[$RANDOM % ${#SURNAMES[@]}]}"
  POSITION="${POSITIONS[$RANDOM % ${#POSITIONS[@]}]}"
  TYPE="${TYPES[$RANDOM % ${#TYPES[@]}]}"
  SALARY=$((3000 + RANDOM % 12000))

  curl -s -X POST http://localhost:8080/employees \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"$NAME\",\"position\":\"$POSITION\",\"salary\":$SALARY,\"type\":\"$TYPE\"}" > /dev/null
done
echo "400 registros inseridos!"
