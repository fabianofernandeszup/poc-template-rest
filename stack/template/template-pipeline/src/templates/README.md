# Templates para pipelines de CI e CD via Gitlab

O objetivo desse repositório é oferecer a Squad de Desenvolvimento templates prontos para realizar Integração Continua (CI) e Deploy (CD) em seus projetos. 

Os templates são formados por conjunto Linguagem + Serviços, mas caso necessário podem ser utilizados de forma separada.

> A sugestão é sempre utilizar um padrão de arquitetura referencia da Rede com um template ja definido.

## CI 

### Tecnologias/Linguagens suportadas:

- angular
- dotnetcore (.net-core)
- java (maven)
- nodejs
- python
- statics (html, css, etc...)
- swagger (json)
- swagger (yaml)

## CD 

### Serviços AWS suportados:

- apigateway
- ecs (docker)
- lambda
- s3
- ec2 (via imagem deploy_puppet)

### Serviços CTMM suportados:

- Puppet Itau (em breve)

## Templates formados:

Os templates são formados por um template de CI (Integração Continua) e um CD (Deploy), por exemplo:

- **ci**: angular | **cd**: s3
- **ci**: swager.json | **cd**: apigateway
- **ci**: swager.yaml | **cd**: apigateway
- **ci**: java maven | **cd**: ecs

### Exemplo

Supondo que você tenha um projeto em ***python*** que precisa ser entregue na Aws ***lambda***, será necessário adicionar ao seu projeto:

- arquivo de template: `python_lambda.yaml`

```
# gitlab-ci.yml
include:
  - project: 'DevOps/templates/ci_cd'
    file: '/python_lambda.yaml'

stages:
  - build                 # jobs: build (branches e tags)
  - test                  # jobs: unit_test (branches e tags)
  - inspection            # jobs: quality_check (branches), quality (tags), security_check (branches), security (tags)
  - package               # jobs: package (branches e tags)
  - publish               # jobs: publish (tags)
  - deploy_dev            # jobs: deploy_dev (master)
  - deploy_hom            # jobs: deploy_hom (tags)
  - governance            # jobs: create_jira_ticket (tags)
  - deploy_prd            # jobs: deploy_prd (tags)
```

Mais detalhes na documentação: [http://devops.pages.gitlab.prd.useredecloud/docs/templates_pipelines/](http://devops.pages.gitlab.prd.useredecloud/docs/templates_pipelines/)
